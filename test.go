package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/tonyd3/propel-gtm/api/clients"
	"github.com/tonyd3/propel-gtm/api/logging"
	"github.com/tonyd3/propel-gtm/api/models"
	"github.com/tonyd3/propel-gtm/api/types"
	"github.com/tonyd3/propel-gtm/api/utils"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// ReviewPullRequest performs a code review on a specific pull request
func (w *CodeReviewWorkflow) ReviewPullRequest(prNumber int, contextBuilder *ContextBuilder, checkIfAlreadyApproved bool, checkExistingComments bool) ([]*InternalReviewComment, []clients.PullRequestFile, error) {
	logger := logging.GetGlobalLogger()
	requestID := "" // Initialize with empty string since we don't have requestID in this context

	// Log workflow start
	logging.LogWorkflowStep(
		"CODE_REVIEW",
		prNumber,
		requestID,
		"review_pull_request_start",
		map[string]interface{}{
			"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			"commit_sha": w.githubConfig.CommitSHA,
		},
	)

	prDetails, err := w.GetPullRequestDetails(prNumber)
	if err != nil {
		logging.LogWorkflowError(
			"CODE_REVIEW",
			prNumber,
			requestID,
			err,
			map[string]interface{}{
				"step":       "get_pr_details",
				"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
				"duration":   time.Since(prDetailsStart),
			},
		)
		return nil, nil, fmt.Errorf("failed to get PR details: %w", err)
	}
	//
	//aiLabels := GetAIReviewLabels(prDetails.Labels)
	//
	//if slices.Contains(aiLabels, "ai:skip-review") {
	//	logging.GetGlobalLogger().Info("AI review skipped due to PR label",
	//		zap.String("reason", "ai:skip-review"),
	//		zap.Int("pr", prDetails.Number),
	//		zap.String("repo", prDetails.RepoName),
	//	)
	//	return nil, nil, fmt.Errorf("AI review skipped due to label 'ai:skip-review'")
	//}

	logging.LogWorkflowStep(
		"CODE_REVIEW",
		prNumber,
		requestID,
		"get_pr_details",
		map[string]interface{}{
			"duration":   time.Since(prDetailsStart),
			"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			"pr_title":   prDetails.Title,
			"commit_sha": prDetails.Head.SHA,
		},
	)
	fmt.Printf("Reviewing PR #%d: %s\n", prNumber, prDetails.Title)

	// Step 3: Double Check if a review is still required, skip review if already approved
	if checkIfAlreadyApproved && w.IsPRAlreadyApproved(prNumber) {
		logging.LogWorkflowStep(
			"CODE_REVIEW",
			prNumber,
			requestID,
			"check_already_approved",
			map[string]interface{}{
				"approved":   true,
				"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			},
		)
		logger.Info("PR has already been approved", zap.Int("prNumber", prNumber))
		// Return nil for files as it's a successful early exit
		return nil, nil, nil
	}

	// Step 4: Get pull request files
	filesStart := time.Now()
	files, err := w.githubConfig.Client.GetPullRequestFiles(
		w.githubConfig.Token,
		w.githubConfig.Owner,
		w.githubConfig.Repo,
		prNumber)
	if err != nil {
		logging.LogWorkflowError(
			"CODE_REVIEW",
			prNumber,
			requestID,
			err,
			map[string]interface{}{
				"step":       "get_pr_files",
				"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
				"duration":   time.Since(filesStart),
			},
		)
		return nil, nil, fmt.Errorf("failed to get PR files: %w", err)
	}
	logging.LogWorkflowStep(
		"CODE_REVIEW",
		prNumber,
		requestID,
		"get_pr_files",
		map[string]interface{}{
			"duration":   time.Since(filesStart),
			"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			"file_count": len(files),
		},
	)
	logging.GetGlobalLogger().Info("Found changed files", zap.Int("count", len(files)))

	// Store files in workflow for later use
	w.prFiles = files

	// Extract filenames for dependency analysis
	filePaths := make([]string, len(files))
	for i, file := range files {
		filePaths[i] = file.Filename
	}

	existingComments := []clients.PullRequestComment{}
	if checkExistingComments {
		// Fetch existing comments for deduplication
		commentsStart := time.Now()
		existingComments, err = w.FetchExistingComments(prNumber)
		existingComments = filterOutExternalBotComments(existingComments)
		if err != nil {
			logger.Warn("Failed to fetch existing comments",
				zap.Error(err),
				zap.Int("pr_number", prNumber),
				zap.String("repository", w.githubConfig.Owner+"/"+w.githubConfig.Repo))
			// Continue with the review even if we can't fetch existing comments
			existingComments = []clients.PullRequestComment{}
		}
		logging.LogWorkflowStep(
			"CODE_REVIEW",
			prNumber,
			requestID,
			"fetch_existing_comments",
			map[string]interface{}{
				"duration":      time.Since(commentsStart),
				"repository":    w.githubConfig.Owner + "/" + w.githubConfig.Repo,
				"comment_count": len(existingComments),
			},
		)
		fmt.Printf("Found %d existing comments\n", len(existingComments))
	}

	// Build unified author context from PR description and comments
	authorContextStart := time.Now()
	authorContextBuilder := NewAuthorContextBuilder(w.githubConfig)
	authorContext, err := authorContextBuilder.BuildAuthorContext(prNumber, prDetails.Body, prDetails.User.Login)
	if err != nil {
		logger.Warn("Failed to build author context",
			zap.Error(err),
			zap.Int("pr_number", prNumber),
			zap.String("repository", w.githubConfig.Owner+"/"+w.githubConfig.Repo),
			zap.String("author", prDetails.User.Login))
		// Continue with the review even if we can't build author context
		authorContext = &AuthorContext{
			HasContent: false,
			PRAuthor:   prDetails.User.Login,
		}
	}
	logging.LogWorkflowStep(
		"CODE_REVIEW",
		prNumber,
		requestID,
		"build_author_context",
		map[string]interface{}{
			"duration":           time.Since(authorContextStart),
			"repository":         w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			"has_content":        authorContext.HasContent,
			"has_description":    len(authorContext.UserDescription) > 0,
			"description_length": len(authorContext.UserDescription),
			"comment_count":      len(authorContext.AuthorComments),
		},
	)
	commit := prDetails.Head.SHA
	fmt.Printf("Latest commit: %s\n", commit)

	// Set the commit SHA and base SHA in the GitHub config
	w.githubConfig.WithCommit(commit, prDetails.Base.SHA)

	// Set up GitHub configuration for the context builder
	contextStart := time.Now()
	if contextBuilder != nil {
		// Share the GitHub config with the context builder
		contextBuilder.WithGitHubConfig(w.githubConfig)

		// Add dependency context to the context builder
		contextBuilder.WithDependencyContext(filePaths)

		// Add language-specific context to the context builder
		contextBuilder.WithLanguageSpecificContext(filePaths)

		// Add code index context to the context builder (with default of 2 degrees)
		// contextBuilder.WithCodeIndexContextLegacy(filePaths, 2)

		// Add the new V2 code index context
		contextBuilder.WithCodeIndexContextGeneric(files, filePaths)

		// Add the knowledge base context (includes file-based rules) if we have a company ID
		if companyId := w.repoWorkflowSetting.CompanyId; companyId > 0 {
			contextBuilder.WithKnowledgeBaseContext(companyId)
			logging.GetGlobalLogger().Info("Added knowledge base context (including file-based rules) for company ID", zap.Uint("company_id", companyId))
		}
	}

	// Step 5: Prepare context message for AI
	var additionalContext map[string]interface{}
	if contextBuilder != nil {
		var err error
		additionalContext, err = contextBuilder.Build()
		if err != nil {
			logger.Warn("Failed to build additional context",
				zap.Error(err),
				zap.Int("pr_number", prNumber),
				zap.String("repository", w.githubConfig.Owner+"/"+w.githubConfig.Repo))
			additionalContext = map[string]interface{}{}
		}
		// Remove analysis results that contain errors
		errorKeywords := []string{"error", "fail", "timeout", "exceeded", "invalid"}
		for _, key := range []string{"code_index_py", "code_index_ts_v2", "code_index"} {
			if analysis, ok := additionalContext[key].(map[string]interface{}); ok {
				hasError := false
				// Check all keys and string values for error keywords
				for k, v := range analysis {
					// Convert key to lowercase for case-insensitive matching
					kLower := strings.ToLower(k)
					for _, keyword := range errorKeywords {
						if strings.Contains(kLower, keyword) {
							hasError = true
							break
						}
						// If value is string, check it too
						if strVal, ok := v.(string); ok {
							if strings.Contains(strings.ToLower(strVal), keyword) {
								hasError = true
								break
							}
						}
					}
					if hasError {
						break
					}
				}
				if hasError {
					log.Printf("\n\nAdditional context: %s\n\n", additionalContext)
					log.Printf("Removing %s from context due to detected error keywords\n", key)
					delete(additionalContext, key)
				}
			}
		}
	} else {
		additionalContext = map[string]interface{}{}
	}

	useSeparateDuplicateDetection := models.IsFeatureEnabledForCompany(w.db, string(types.SeparateDuplicateDetection), w.repoWorkflowSetting.CompanyId)
	// If separate duplicate detection is DISABLED, include existing comments in the main AI context (original behavior)
	if !useSeparateDuplicateDetection && len(existingComments) > 0 {
		additionalContext["existing_comments"] = existingComments
		additionalContext[string(types.SeparateDuplicateDetection)] = true
		logger.Info("[SeparateDuplicateDetection] Feature flag disabled - including existing comments in main AI context",
			zap.Int("existing_comments_count", len(existingComments)),
			zap.Int("pr_number", prNumber))
	} else if useSeparateDuplicateDetection {
		logger.Info("[SeparateDuplicateDetection] Feature flag enabled - will use separate duplicate detection service",
			zap.Int("existing_comments_count", len(existingComments)),
			zap.Int("pr_number", prNumber))
	}

	// Add unified author context
	if authorContext.HasContent {
		additionalContext["author_context"] = authorContext
	}

	// Add commitable suggestions flag from configuration (default to false for backward compatibility)
	if w.config != nil {
		additionalContext["committable_suggestions_enabled"] = w.config.CommittableSuggestions
	} else {
		additionalContext["committable_suggestions_enabled"] = true
	}
	additionalContext["workflow_company_id"] = w.repoWorkflowSetting.CompanyId

	logging.LogWorkflowStep(
		"CODE_REVIEW",
		prNumber,
		requestID,
		"prepare_context",
		map[string]interface{}{
			"duration":              time.Since(contextStart),
			"repository":            w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			"context_keys":          getContextKeys(additionalContext),
			"has_existing_comments": len(existingComments) > 0,
		},
	)

	// Prepare context and user messages
	messageStart := time.Now()

	// Check if persona-based review is enabled
	usePersonas := models.IsFeatureEnabledForCompany(w.db, string(types.CodeReviewPersonas), w.repoWorkflowSetting.CompanyId)
	var contextMessage string

	if usePersonas && w.config != nil && len(w.config.ActivePersonas) > 0 {
		// Use persona-enhanced system message
		personaService := NewCodeReviewPersonaService()
		contextMessage = personaService.BuildPersonaSystemMessage(additionalContext, w.config.ActivePersonas, commit)
		logging.GetGlobalLogger().Info("Using persona-based code review",
			zap.Int("pr_number", prNumber),
			zap.Int("active_personas", len(w.config.ActivePersonas)),
			zap.Any("personas", w.config.ActivePersonas))
	} else {
		// Use default system message
		contextMessage = buildSystemMessage(additionalContext, commit)
	}

	contextMessageTokenCount := CountTokens(contextMessage)

	// Step 6: Prepare user message with file changes - assuming we take up 75% of the token budget to be conservative
	tokenBudget := (MaxAllowedTokens - contextMessageTokenCount) * 75 / 100
	userMessage := w.prepareUserMessage(files, tokenBudget)

	if models.IsFeatureEnabledForCompany(w.db, string(types.AgenticWorkflow), w.repoWorkflowSetting.CompanyId) {
		ragContext := w.prepareContextFromRAG(files, fmt.Sprintf("%s-%s", w.githubConfig.Owner, w.githubConfig.Repo))
		logger.Info("RAG context", zap.String("rag_context", ragContext))
	}

	// Log context message to file for debugging
	if err := logContextMessage(contextMessage, userMessage, commit, prNumber, w, files, contextMessageTokenCount); err != nil {
		logging.LogWorkflowError(
			"CODE_REVIEW",
			prNumber,
			requestID,
			err,
			map[string]interface{}{
				"step":       "log_context_message",
				"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			},
		)
		// Return files because we have them, but also the error
		return nil, files, err
	}

	userMessageTokenCount := CountTokens(userMessage)
	logger.Debug("Token count for userMessage",
		zap.String("commit", commit),
		zap.Int("token_count", userMessageTokenCount))

	combinedMessage := fmt.Sprintf("%s\n%s", contextMessage, userMessage)
	tokenCount := CountTokens(combinedMessage)
	fmt.Printf("Size of the current combinedMessage tokens: %d\n", tokenCount)

	contextMessage, userMessage = PruneToTokenLimit(contextMessage, userMessage, additionalContext, commit, files, MaxAllowedTokens)
	combinedMessage = fmt.Sprintf("%s\n%s", contextMessage, userMessage)
	tokenCount = CountTokens(combinedMessage)
	fmt.Printf("New Size of the current combinedMessage tokens: %d\n", tokenCount)

	logging.LogWorkflowStep(
		"CODE_REVIEW",
		prNumber,
		requestID,
		"prepare_messages",
		map[string]interface{}{
			"duration":            time.Since(messageStart),
			"repository":          w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			"context_token_count": contextMessageTokenCount,
			"user_token_count":    userMessageTokenCount,
			"total_token_count":   tokenCount,
		},
	)

	// Step 7: Generate AI review
	aiStart := time.Now()
	logging.LogWorkflowStep(
		"CODE_REVIEW",
		prNumber,
		requestID,
		"call_ai_model_start",
		map[string]interface{}{
			"repository":  w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			"token_count": tokenCount,
		},
	)

	// Call all three AI models and merge their comments together
	internalComments, err := w.callMultipleAIModels(
		contextMessage,
		userMessage,
		tokenCount,
		additionalContext,
		commit,
		files,
		prNumber,
		requestID,
		aiStart,
	)
	if err != nil {
		logging.GetGlobalLogger().Warn("Failed to generate AI review", zap.Error(err))
		// Return files as we have them, but also the error from AI model call
		return nil, files, err
	}

	// NEW: Apply separate duplicate detection service AFTER AI generation
	if useSeparateDuplicateDetection {
		duplicateDetectionStart := time.Now()
		duplicateDetectionService := NewDuplicateDetectionService(w.aiConfig)
		internalComments, err = duplicateDetectionService.FilterDuplicates(internalComments, existingComments, prNumber, requestID)
		if err != nil {
			logger.Warn("Failed to apply duplicate detection",
				zap.Error(err),
				zap.Int("pr_number", prNumber),
				zap.String("repository", w.githubConfig.Owner+"/"+w.githubConfig.Repo))
			// Continue with original comments if duplicate detection fails
		}
		logging.LogWorkflowStep(
			"CODE_REVIEW",
			prNumber,
			requestID,
			"duplicate_detection",
			map[string]interface{}{
				"duration":   time.Since(duplicateDetectionStart),
				"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			},
		)
	}

	// Validate Review Comments
	validateReviews := models.IsFeatureEnabledForCompany(w.db, string(types.ValidateReviews), w.repoWorkflowSetting.CompanyId)
	if validateReviews {
		internalComments, err = w.ReviewComments(internalComments, files, authorContext, prNumber, requestID, w.aiConfig.CallOpenAI, w.aiConfig.GetOpenAIModel())
		if err != nil {
			logging.GetGlobalLogger().Warn("Failed to validate comments", zap.Error(err))
		}
	}

	validateReviewsGemini := models.IsFeatureEnabledForCompany(w.db, string(types.ValidateReviewsGemini), w.repoWorkflowSetting.CompanyId)
	if validateReviewsGemini {
		internalComments, err = w.ReviewComments(internalComments, files, authorContext, prNumber, requestID, w.aiConfig.CallGemini, w.aiConfig.GetGeminiModel())
		if err != nil {
			logging.GetGlobalLogger().Warn("Failed to validate comments", zap.Error(err))
		}
	}

	validateReviewsOpus := models.IsFeatureEnabledForCompany(w.db, string(types.ValidateReviewsOpus), w.repoWorkflowSetting.CompanyId)
	if validateReviewsOpus {
		internalComments, err = w.ReviewComments(internalComments, files, authorContext, prNumber, requestID, func(contextMessage, userMessage string) (string, error) {
			return w.aiConfig.CallAnthropicWithModel(contextMessage, userMessage, anthropic.ModelClaudeOpus4_20250514)
		}, string(w.aiConfig.GetAnthropicModel()))

		if err != nil {
			logging.GetGlobalLogger().Warn("Failed to validate comments", zap.Error(err))
		}
	}

	// Get Previously Provided Comments
	previousComments, err := w.githubConfig.Client.GetPullRequestReviewComments(w.githubConfig.Token, w.githubConfig.Owner, w.githubConfig.Repo, prNumber)
	if err != nil {
		logging.GetGlobalLogger().Error("Failed to get PR comments", zap.Error(err))
	}
	// Remove comments left by external bots so we don't mis‑detect duplicates
	previousComments = filterOutExternalBotComments(previousComments)
	// Use a token budget of 0 as we only need patch structure, not full file content for this validation.
	parsedFileChanges := ExtractAllFileChanges(w.githubConfig, files, 0)

	for _, comment := range internalComments {
		if len(comment.RejectionReason) != 0 {
			continue
		}

		enableNoOpValidation := models.IsFeatureEnabledForCompany(w.db, string(types.NoOpSuggestionValidation), w.repoWorkflowSetting.CompanyId)
		if enableNoOpValidation && w.suggestionValidator.IsSuggestionNoOp(*comment, files) {
			comment.RejectionReason = "No-op suggestion: suggested code is the same as original code"
			rejectionModel := "no_op_detection"
			comment.RejectionModel = &rejectionModel
			logging.GetGlobalLogger().Info("Skipping no-op suggestion comment",
				zap.Int("pr_number", prNumber),
				zap.String("path", comment.Path),
				zap.Int("line", comment.Line),
				zap.String("body", comment.Body))
			continue
		}

		isAppropriate, reason, err := w.ValidateComment(*comment, previousComments, parsedFileChanges)
		if err != nil {
			logging.GetGlobalLogger().Warn("Failed to validate comment", zap.Error(err))
		} else {
			if !isAppropriate {
				comment.RejectionReason = reason
				currentModel := w.aiConfig.GetOpenAIModel()
				comment.RejectionModel = &currentModel
			}
		}
	}

	// Add commitable suggestion header after all validation is complete
	for _, comment := range internalComments {
		if len(comment.RejectionReason) == 0 && strings.Contains(comment.Body, "```suggestion") {
			comment.Body += "\n\n⚡ **Committable suggestion**\n\n" +
				"Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, " +
				"contains no missing lines, and has no issues with indentation."
		}
	}

	// Return successfully generated comments and the files
	return internalComments, files, nil
}

type CommentResponse struct {
	Id          string       `json:"id"`
	Valid       StringOrBool `json:"valid"`
	Explanation string       `json:"explanation"`
}

type StringOrBool string

func (sb *StringOrBool) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*sb = StringOrBool(s)
		return nil
	}

	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			*sb = "true"
		} else {
			*sb = "false"
		}
		return nil
	}

	return json.Unmarshal(data, (*string)(sb))
}

func (sb StringOrBool) IsTrue() bool {
	lower := strings.ToLower(strings.TrimSpace(string(sb)))
	return lower == "true" || lower == "yes"
}

func BuildFullPatch(files []clients.PullRequestFile) string {
	var sb strings.Builder
	for _, file := range files {
		sb.WriteString(fmt.Sprintf("--- %s\n+++ %s\n%s\n\n", file.Filename, file.Filename, file.Patch))
	}
	return sb.String()
}

// callMultipleAIModels calls all three AI models in parallel using goroutines and merges their comments together
func (w *CodeReviewWorkflow) callMultipleAIModels(
	contextMessage string,
	userMessage string,
	tokenCount int,
	additionalContext map[string]interface{},
	commit string,
	files []clients.PullRequestFile,
	prNumber int,
	requestID string,
	aiStart time.Time,
) ([]*InternalReviewComment, error) {
	// Create channels for results and errors
	type modelResult struct {
		comments  []*InternalReviewComment
		err       error
		modelName string
	}

	useGemini := models.IsFeatureEnabledForCompany(w.db, string(types.AddGeminiResults), w.repoWorkflowSetting.CompanyId)
	useOpenAI := models.IsFeatureEnabledForCompany(w.db, string(types.AddOpenAIResults), w.repoWorkflowSetting.CompanyId)

	// Determine how many models we'll call based on feature flags
	modelCount := 1 // Anthropic is always enabled
	if useOpenAI {
		modelCount++
	}
	if useGemini {
		modelCount++
	}

	resultChan := make(chan modelResult, modelCount)
	go func() {
		comments, err := w.callAIModel(
			contextMessage,
			userMessage,
			tokenCount,
			additionalContext,
			commit,
			files,
			prNumber,
			requestID+"-anthropic",
			aiStart,
			w.aiConfig.CallAnthropic,
			"anthropic",
			string(w.aiConfig.GetAnthropicModel()),
		)
		resultChan <- modelResult{comments, err, "Anthropic"}
	}()

	if useOpenAI {
		go func() {
			comments, err := w.callAIModel(
				contextMessage,
				userMessage,
				tokenCount,
				additionalContext,
				commit,
				files,
				prNumber,
				requestID+"-openai",
				aiStart,
				w.aiConfig.CallOpenAI,
				"openai",
				w.aiConfig.GetOpenAIModel(),
			)
			resultChan <- modelResult{comments, err, "OpenAI"}
		}()
	}

	if useGemini {
		go func() {
			comments, err := w.callAIModel(
				contextMessage,
				userMessage,
				tokenCount,
				additionalContext,
				commit,
				files,
				prNumber,
				requestID+"-gemini",
				aiStart,
				w.aiConfig.CallGemini,
				"google",
				w.aiConfig.GetGeminiModel(),
			)
			resultChan <- modelResult{comments, err, "Gemini"}
		}()

	}

	// Collect results
	var anthropicComments, openaiComments, geminiComments []*InternalReviewComment
	var err1, err2, err3 error

	// Wait for all results based on the number of models we're calling
	for range modelCount {
		result := <-resultChan
		switch result.modelName {
		case "Anthropic":
			anthropicComments = result.comments
			err1 = result.err
			if err1 != nil {
				logging.GetGlobalLogger().Error("Error calling Anthropic", zap.Error(err1))
			}
		case "OpenAI":
			openaiComments = result.comments
			err2 = result.err
			if err2 != nil {
				logging.GetGlobalLogger().Error("Error calling OpenAI", zap.Error(err2))
			}
		case "Gemini":
			geminiComments = result.comments
			err3 = result.err
			if err3 != nil {
				logging.GetGlobalLogger().Error("Error calling Gemini", zap.Error(err3))
			}
		}
	}

	// Check if all calls failed
	if err1 != nil && err2 != nil && err3 != nil {
		return nil, fmt.Errorf("all AI model calls failed: Anthropic: %v, OpenAI: %v, Gemini: %v", err1, err2, err3)
	}

	// Merge all comments, ensuring uniqueness
	var mergedComments []*InternalReviewComment

	// Helper function to check if a comment overlaps with any existing comments
	isOverlapping := func(newComment *InternalReviewComment, existingComments []*InternalReviewComment) bool {
		for _, existing := range existingComments {
			if commentsOverlap(newComment.PullRequestComment, existing.PullRequestComment) {
				return true
			}
		}
		return false
	}

	// Add Anthropic comments first (prioritize these)
	mergedComments = append(mergedComments, anthropicComments...)

	// Add OpenAI comments if they don't overlap with existing comments
	for _, comment := range openaiComments {
		if comment.Body != "" && !isOverlapping(comment, mergedComments) {
			mergedComments = append(mergedComments, comment)
		}
	}

	// Add Gemini comments if they don't overlap with existing comments
	for _, comment := range geminiComments {
		if comment.Body != "" && !isOverlapping(comment, mergedComments) {
			mergedComments = append(mergedComments, comment)
		}
	}

	logging.LogWorkflowStep(
		"CODE_REVIEW",
		prNumber,
		requestID,
		"merged_ai_models_complete",
		map[string]interface{}{
			"repository":      w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			"duration":        time.Since(aiStart),
			"anthropic_count": len(anthropicComments),
			"openai_count":    len(openaiComments),
			"openai_enabled":  useOpenAI,
			"gemini_count":    len(geminiComments),
			"gemini_enabled":  useGemini,
			"merged_count":    len(mergedComments),
		},
	)

	for _, comment := range mergedComments {
		commentType, err := w.classifyCommentType(
			comment.Body,
			BuildFullPatch(files),
			commit,
			requestID+"-classify-"+comment.ID,
		)
		if err != nil {
			logging.GetGlobalLogger().Warn("Failed to classify type", zap.Error(err))
			continue
		}
		comment.Type = commentType
		logging.GetGlobalLogger().Info("Comment type classified successfully", zap.String("comment_type", commentType), zap.String("comment_body", comment.Body))
	}

	titleCaser := cases.Title(language.English)
	for _, comment := range mergedComments {
		if comment.Type != "" {
			titleCasedType := strings.ReplaceAll(titleCaser.String(strings.ReplaceAll(comment.Type, "_", " ")), " ", "")
			formattedType := "[**" + titleCasedType + "**]"
			prefixWithBrackets := "[" + titleCasedType + "]"

			body := comment.Body
			// Remove any existing prefix, formatted or simple.
			if strings.HasPrefix(body, formattedType) {
				body = strings.TrimPrefix(body, formattedType)
			} else if strings.HasPrefix(body, prefixWithBrackets) {
				body = strings.TrimPrefix(body, prefixWithBrackets)
			}

			// Clean up any leading space or newlines from the actual content.
			body = strings.TrimLeft(body, " \n")

			// Re-apply the correctly formatted prefix.
			comment.Body = formattedType + "\n\n" + body
		}
	}

	return mergedComments, nil
}

// callAIModel handles calling an AI model with appropriate error handling and retry logic
func (w *CodeReviewWorkflow) callAIModel(
	contextMessage string,
	userMessage string,
	tokenCount int,
	additionalContext map[string]interface{},
	commit string,
	files []clients.PullRequestFile,
	prNumber int,
	requestID string,
	aiStart time.Time,
	modelCallFn func(string, string) (string, error),
	provider string,
	model string,
) ([]*InternalReviewComment, error) {
	message, err := modelCallFn(contextMessage, userMessage)

	if err != nil {
		currentTokenCount, errMsg := extractTokenCountFromAnthropicError(err)
		if errMsg == nil {
			adjustedMaxTokenLimit := computeAdjustedTokenLimit(tokenCount, currentTokenCount, 3.0)
			SendSlackNotification(1, fmt.Sprintf(
				"🤖 Anthropic token mismatch detected.\n• Count we estimated: %d\n• Anthropic actual: %d\n• Adjusted Max Limit: %d",
				tokenCount, currentTokenCount, adjustedMaxTokenLimit,
			))

			logging.LogWorkflowError(
				"CODE_REVIEW",
				prNumber,
				requestID,
				err,
				map[string]interface{}{
					"step":                  "call_ai_model",
					"repository":            w.githubConfig.Owner + "/" + w.githubConfig.Repo,
					"duration":              time.Since(aiStart),
					"error_type":            "token_limit",
					"estimated_token_count": tokenCount,
					"actual_token_count":    currentTokenCount,
					"adjusted_token_limit":  adjustedMaxTokenLimit,
				},
			)

			// Log the retry attempt
			retryStart := time.Now()
			logging.LogWorkflowStep(
				"CODE_REVIEW",
				prNumber,
				requestID,
				"retry_ai_model_call",
				map[string]interface{}{
					"repository":           w.githubConfig.Owner + "/" + w.githubConfig.Repo,
					"adjusted_token_limit": adjustedMaxTokenLimit,
				},
			)

			contextMessage, userMessage = PruneToTokenLimit(contextMessage, userMessage, additionalContext, commit, files, adjustedMaxTokenLimit)

			message, err = modelCallFn(contextMessage, userMessage)

			if err != nil {
				SendSlackNotification(1, fmt.Sprintf("❌ Retried after adjustment but still failed: %v", err))
				logging.LogWorkflowError(
					"CODE_REVIEW",
					prNumber,
					requestID,
					err,
					map[string]interface{}{
						"step":       "retry_ai_model_call",
						"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
						"duration":   time.Since(retryStart),
						"error_type": "retry_failed",
					},
				)
			} else {
				SendSlackNotification(1, "✅ Retried after adjustment and succeeded.")
				logging.LogWorkflowStep(
					"CODE_REVIEW",
					prNumber,
					requestID,
					"retry_ai_model_success",
					map[string]interface{}{
						"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
						"duration":   time.Since(retryStart),
					},
				)
			}
		} else {
			SendSlackNotification(1, fmt.Sprintf("⚠️ Failed to extract token count from error: %v", errMsg))
			logging.LogWorkflowError(
				"CODE_REVIEW",
				prNumber,
				requestID,
				err,
				map[string]interface{}{
					"step":       "call_ai_model",
					"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
					"duration":   time.Since(aiStart),
					"error_type": "unknown",
				},
			)
		}
	} else {
		// Log successful AI model call
		responseTokens := logging.EstimateTokenCount(message)
		logging.LogModelCall(
			"CODE_REVIEW",
			prNumber,
			requestID,
			"anthropic",
			"code_review",
			tokenCount,
			responseTokens,
			time.Since(aiStart),
		)

		logging.LogWorkflowStep(
			"CODE_REVIEW",
			prNumber,
			requestID,
			"call_ai_model_complete",
			map[string]interface{}{
				"repository":      w.githubConfig.Owner + "/" + w.githubConfig.Repo,
				"duration":        time.Since(aiStart),
				"response_tokens": responseTokens,
			},
		)
	}

	var internalComments []*InternalReviewComment
	parseStart := time.Now()
	err = SanitizeAndParseJSON(message, &internalComments)
	if err != nil {
		logging.LogWorkflowError(
			"CODE_REVIEW",
			prNumber,
			requestID,
			err,
			map[string]interface{}{
				"step":       "parse_comments",
				"repository": w.githubConfig.Owner + "/" + w.githubConfig.Repo,
				"duration":   time.Since(parseStart),
			},
		)
		fmt.Printf("Failed to parse comments: %v\n", err)
		return nil, fmt.Errorf("failed to parse comments: %w", err)
	}

	for _, c := range internalComments {
		c.Provider = provider
		c.Model = model
		c.CommitID = commit
	}

	// Log workflow completion
	logging.LogWorkflowComplete(
		"CODE_REVIEW",
		prNumber,
		requestID,
		aiStart,
		map[string]interface{}{
			"repository":     w.githubConfig.Owner + "/" + w.githubConfig.Repo,
			"comments_count": len(internalComments),
			"commit_sha":     commit,
			"provider":       provider,
			"model":          model,
		},
	)

	return internalComments, nil
}

/*
var comments []*clients.PullRequestComment
*/
func (w *CodeReviewWorkflow) AddExecutionLog(logMessage string) {
	w.codeWorkflowExecution.AddExecutionLog(w.db, models.WorkflowTypeCODE_REVIEW, logMessage)
}

// MaxAllowedTokens defines the maximum number of tokens allowed in the prompt.
const MaxAllowedTokens = 200000

// buildSystemMessage constructs the system message JSON as a string using the modular system message builder.
func buildSystemMessage(additionalContext map[string]interface{}, commit string) string {
	companyId, _ := additionalContext["workflow_company_id"].(uint)

	// Create system message builder
	builder := NewSystemMessageBuilder(commit, companyId)

	// Define base configuration for standard code review
	config := SystemMessageConfig{
		Role:           "You are a world class software engineer and an expert in code review.",
		Objective:      "You are conducting a code review for another member of your team. Provide ONLY specific, actionable, and concise feedback that directly improves code quality.",
		Guidelines:     "Focus exclusively on substantive issues. If you have any feedback, provide code snippets or specific suggestions with examples in Markdown format.",
		ThoughtProcess: "Think deeply and reason about how a world-class engineer would approach this code review. Think through all possibilities and trade-offs, critique them, refine your thinking, and then focus on only feedback that is actionable and relevant. IMPORTANT: Do not comment just to acknowledge that code is already correct or follows best practices. Only provide comments when there is a concrete improvement or correction to suggest.",
		Brevity:        "Be concise and focused in your review, do not include any reviews that might be subjective or are not actionable. Having extra comments that are not actionable will not improve code quality, it will go against the best practices of code review.",
	}

	// Build base message
	message := builder.BuildBaseMessage(config, additionalContext)

	// Add modular components
	builder.AddCommittableSuggestions(message, additionalContext)
	builder.AddDuplicateDetection(message, additionalContext)
	builder.AddKnowledgeBaseGuidelines(message, additionalContext)

	// Convert to JSON
	return builder.ToJSON(message)
}

// prepareUserMessage creates the user message with file changes
func (w *CodeReviewWorkflow) prepareUserMessage(files []clients.PullRequestFile, tokenBudget int) string {
	fileChanges := ExtractAllFileChanges(w.githubConfig, files, tokenBudget)

	message, err := json.Marshal(fileChanges)

	if err != nil || len(fileChanges) == 0 {
		logging.GetGlobalLogger().Error("Error marshaling file changes", zap.Error(err))
		// Resort to basic encoding
		messages := []string{}
		for _, file := range files {
			messages = append(messages, fmt.Sprintf(`{
				 "path": "%s",
				 "additions": %d,
				 "deletions": %d,
				 "changes": %d,
				 "status": "%s",
				 "patch": "%s"
			 }`, file.Filename, file.Additions, file.Deletions, file.Changes, file.Status, file.Patch))
		}
		return fmt.Sprintf(`{"file_changes": [%s]}`, strings.Join(messages, ","))
	}

	return fmt.Sprintf(`{"file_changes": %s}`, string(message))
}

func (w *CodeReviewWorkflow) prepareContextFromRAG(files []clients.PullRequestFile, collectionName string) string {
	queries, err := GenerateQueriesFromPRFiles(w.githubConfig, w.aiConfig, files)
	if err != nil {
		logging.GetGlobalLogger().Error("Error generating queries", zap.Error(err))
		return ""
	}

	chromaOutput, err := ExecuteChromaSearch(queries, collectionName)
	if err != nil {
		logging.GetGlobalLogger().Error("Error executing chroma search", zap.Error(err))
		return ""
	}

	return chromaOutput
}

type FilteredPullRequestComment struct {
	clients.PullRequestComment
	Reason string `json:"reason"`
}

// PostReviewComments posts the review comments to GitHub after validating them
func (w *CodeReviewWorkflow) PostReviewComments(prNumber int, companyId uint, comments []*InternalReviewComment, pullRequestFiles []clients.PullRequestFile) (postedComments []*InternalReviewComment, filteredComments []*InternalReviewComment, err error) {
	prospectiveComments := []*InternalReviewComment{}
	for _, comment := range comments {
		if len(comment.RejectionReason) > 0 {
			filteredComments = append(filteredComments, comment)
			logging.GetGlobalLogger().Info("Skipping inappropriate comment",
				zap.Int("pr_number", prNumber),
				zap.String("reason", comment.RejectionReason),
				zap.String("body", comment.Body))
			continue
		} else if len(comment.AcceptanceReason) > 0 {
			prospectiveComments = append(prospectiveComments, comment)
			logging.GetGlobalLogger().Info("Accepting comment",
				zap.Int("pr_number", prNumber),
				zap.String("reason", comment.AcceptanceReason),
				zap.String("body", comment.Body))
			continue
		}
	}

	// Apply tiered filtering using the new centralized function
	prospectiveComments = ApplyTieredCommentFiltering(
		w.db,
		companyId,
		logging.GetGlobalLogger(),
		prospectiveComments,
		w.config,
		w.aiConfig,
	)

	// Post the validated comments
	for _, comment := range prospectiveComments {
		response, err := w.githubConfig.Client.PostPullRequestComment(
			w.githubConfig.Token,
			w.githubConfig.Owner,
			w.githubConfig.Repo,
			prNumber,
			&comment.PullRequestComment)

		if err != nil {
			logging.GetGlobalLogger().Error("Failed to post comment", zap.Error(err))
			comment.RejectionReason = "failed to post comment, error: " + err.Error()
			filteredComments = append(filteredComments, comment)
		} else {
			// Record PRComment in database
			prComment := models.PRComment{
				SingleCompanyModel: models.SingleCompanyModel{
					CompanyId: companyId,
				},
				PRNumber:         prNumber,
				Owner:            w.githubConfig.Owner,
				Repo:             w.githubConfig.Repo,
				CommentID:        response.ID,
				Author:           response.User.Login,
				Body:             comment.Body,
				CommentCreatedAt: utils.ParseTimeOrDefault(response.CreatedAt, time.Now()),
				CommentUpdatedAt: utils.ParseTimeOrDefault(response.UpdatedAt, time.Now()),
				IsByWorkflow:     true,
				CommentType:      models.ReviewCommentType(comment.Type),
			}
			if err := w.db.Create(&prComment).Error; err != nil {
				// Log the error but continue with the next comment
				logging.GetGlobalLogger().Error("Failed to record PRComment", zap.Error(err))
			}

			postedComments = append(postedComments, comment)
		}
	}

	// Fall back to the original automatic approval logic if no decision was made
	if w.config.AutomaticApproval && len(postedComments) == 0 {
		// Determine if the PR should be approved based on the comments
		should, reason, err := w.shouldApprovePR(comments)
		if err != nil {
			logging.GetGlobalLogger().Info("Failed to determine if PR should be approved", zap.Error(err))
			// Continue with the default approval logic
		} else if should {
			// Approve the PR with the provided reason
			approveMessage := fmt.Sprintf(`LGTM, ship it! :ship:
			<details>
			<summary>Why was this auto-approved?</summary>
			%s
			</details>`, reason)
			err := w.githubConfig.Client.ApprovePullRequest(w.githubConfig.Token, w.githubConfig.Owner, w.githubConfig.Repo, prNumber, approveMessage)
			if err != nil {
				return postedComments, filteredComments, fmt.Errorf("failed to approve PR: %w", err)
			}
			SendSlackNotification(1, fmt.Sprintf("✅ PR %d approved", prNumber))
			w.AddExecutionLog("PR is auto approved")
			logging.GetGlobalLogger().Info("PR approved", zap.Int("pr_number", prNumber), zap.String("reason", reason))
		} else if reason != "" {
			// Post a comment explaining why the PR was not approved
			logging.GetGlobalLogger().Info("PR not approved", zap.Int("pr_number", prNumber), zap.String("reason", reason))
		}
	}

	return postedComments, filteredComments, nil
}

// shouldApprovePR determines if a PR should be approved based on the comments
func (w *CodeReviewWorkflow) shouldApprovePR(comments []*InternalReviewComment) (bool, string, error) {
	// If there are no comments, we can approve the PR
	if len(comments) == 0 {
		return true, "AI analysis completed with no actionable comments or suggestions.", nil
	}

	// Build a system message for PR approval decision
	systemMessage := `You are a code review approval decision maker. Your task is to determine if a pull request should be approved based on the code review comments.

You should consider:
1. The severity of the issues identified in the comments
2. Whether the issues are blockers or just suggestions
3. The overall quality of the code based on the comments

You should approve if the overall quality is good and there are no major issues.
You should reject if there are major issues that need to be addressed before the PR should be approved.

Respond with ONLY one of these formats:
- "approve: [reason]" if the PR should be approved
- "reject: [reason]" if the PR should not be approved

Reasoning guidelines:
- Your reason must be specific, concise, and grounded in the actual review comments.
- Do not hallucinate or infer anything beyond what is explicitly stated.
- If the review comments contain only vague approvals like "LGTM" or "No issues", you may omit the reason entirely by responding with just: "approve" (no colon, no explanation).
- Avoid repeating "LGTM", "Looks good", or similar phrases in the reason unless you're quoting a reviewer directly for traceability.

Only provide a reason when it adds clarity about **why** the PR is safe to approve based on actual reviewer feedback.`

	// Build the user message with all the comments
	userMessage := "Please evaluate if this pull request should be approved based on the following code review comments:\n\n"
	for i, comment := range comments {
		// Extract comment type from the body if available, otherwise use "general"
		commentType := "general"

		// Try to parse the comment body as JSON to extract the type if it exists
		var commentData map[string]interface{}
		if err := json.Unmarshal([]byte(comment.Body), &commentData); err == nil {
			if typeVal, ok := commentData["type"].(string); ok {
				commentType = typeVal
			}
		}

		userMessage += fmt.Sprintf("Comment %d (type: %s):\n%s\n\n", i+1, commentType, comment.Body)
	}

	// Call the AI to make the approval decision
	response, err := w.aiConfig.CallAnthropic(systemMessage, userMessage)
	if err != nil {
		return false, "", fmt.Errorf("failed to determine if PR should be approved: %w", err)
	}

	// Parse the response
	response = strings.TrimSpace(response)
	if strings.HasPrefix(strings.ToLower(response), "approve:") {
		reason := strings.TrimPrefix(response, "approve:")
		reason = strings.TrimPrefix(reason, "Approve:")
		reason = strings.TrimSpace(reason)
		return true, reason, nil
	} else if strings.HasPrefix(strings.ToLower(response), "reject:") {
		reason := strings.TrimPrefix(response, "reject:")
		reason = strings.TrimPrefix(reason, "Reject:")
		reason = strings.TrimSpace(reason)
		return false, reason, nil
	}

	// If we can't parse the response, return an error
	return false, "", fmt.Errorf("failed to parse approval decision: %s", response)
}

// logContextMessage logs the context message and user message to files for debugging
func logContextMessage(contextMessage, userMessage, commit string, prNumber int, w *CodeReviewWorkflow, files []clients.PullRequestFile, contextMessageTokenCount int) error {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Warning: Failed to create logs directory: %v", err)
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Write new log file with timestamp and metadata
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	metadata := map[string]interface{}{
		"timestamp":   timestamp,
		"commit":      commit,
		"pr_number":   prNumber,
		"owner":       w.githubConfig.Owner,
		"repo":        w.githubConfig.Repo,
		"num_files":   len(files),
		"token_count": contextMessageTokenCount,
	}

	metadataJSON, _ := json.MarshalIndent(metadata, "", "  ")
	logContent := fmt.Sprintf("# Code Review Log\n\n## Metadata\n%s\n\n## Context Message\n%s\n\n## User Message\n%s\n",
		string(metadataJSON), contextMessage, userMessage)

	logFile := filepath.Join(logDir, "context_message.log")

	// Ensure file exists or create it
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		_, err = os.Create(logFile)
		if err != nil {
			log.Printf("Warning: Failed to create log file: %v", err)
			return fmt.Errorf("failed to create log file: %w", err)
		}
	}

	// Write the content
	if err := os.WriteFile(logFile, []byte(logContent), 0644); err != nil {
		log.Printf("Warning: Failed to write context message to log file: %v", err)
		return fmt.Errorf("failed to write context message to log file: %w", err)
	}

	log.Printf("Context message logged to: %s", logFile)

	// Also write a timestamped backup
	backupFile := filepath.Join(logDir, fmt.Sprintf("context_message_%s.log",
		time.Now().Format("20060102_150405")))
	if err := os.WriteFile(backupFile, []byte(logContent), 0644); err != nil {
		log.Printf("Warning: Failed to write backup log file: %v", err)
		// Continue even if backup fails
	} else {
		log.Printf("Backup log written to: %s", backupFile)
	}

	return nil
}

func (w *CodeReviewWorkflow) GetPullRequestDetails(prNumber int) (*clients.PullRequestDetails, error) {
	return w.githubConfig.Client.GetPullRequestDetails(
		w.githubConfig.Token,
		w.githubConfig.Owner,
		w.githubConfig.Repo,
		prNumber)
}

func (w *CodeReviewWorkflow) IsPRAlreadyApproved(prNumber int) bool {
	reviews, err := w.githubConfig.Client.GetPullRequestReviews(
		w.githubConfig.Token,
		w.githubConfig.Owner,
		w.githubConfig.Repo,
		prNumber)
	if err != nil {
		logging.GetGlobalLogger().Error("Failed to get PR reviews", zap.Error(err))
		// Continue with the review even if we can't get previous reviews
	}

	for _, review := range reviews {
		if *review.State == "APPROVED" {
			return true
		}
	}

	return false
}

// getContextKeys extracts the keys from a map for logging purposes
func getContextKeys(contextMap map[string]interface{}) []string {
	keys := make([]string, 0, len(contextMap))
	for k := range contextMap {
		keys = append(keys, k)
	}
	return keys
}
