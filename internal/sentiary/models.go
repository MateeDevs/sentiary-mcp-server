package sentiary

type ProjectInput struct {
	ProjectID string `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
}

type ListStringsInput struct {
	ProjectID          string   `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	Size               int      `json:"size,omitempty" jsonschema:"Maximum number of strings to return. Defaults to the API default."`
	Cursor             string   `json:"cursor,omitempty" jsonschema:"Pagination cursor returned by a previous response."`
	IncludeLanguageIDs []string `json:"includeLanguageIds,omitempty" jsonschema:"Translation language IDs to include in each string."`
	Filter             string   `json:"filter,omitempty" jsonschema:"Optional API filter: ALL, TRANSLATED, UNTRANSLATED, or PARTIALLY_TRANSLATED."`
	FilterLanguage     string   `json:"filterLanguage,omitempty" jsonschema:"Language ID used by filter."`
	Order              string   `json:"order,omitempty" jsonschema:"Optional API order, for example UPDATED_ASC or UPDATED_DESC."`
}

type SearchStringInput struct {
	ProjectID          string   `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	Name               string   `json:"name" jsonschema:"String name to search for. Results are filtered by term name only."`
	Exact              bool     `json:"exact,omitempty" jsonschema:"When true, returns only exact name matches."`
	Size               int      `json:"size,omitempty" jsonschema:"Maximum API page size to search. Defaults to 25."`
	IncludeLanguageIDs []string `json:"includeLanguageIds,omitempty" jsonschema:"Translation language IDs to include in each string."`
}

type GetStringInput struct {
	ProjectID string `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	StringID  string `json:"stringId" jsonschema:"String/term ID."`
}

type GetStringByKeyInput struct {
	ProjectID          string   `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	Key                string   `json:"key" jsonschema:"String key/name to fetch exactly."`
	IncludeLanguageIDs []string `json:"includeLanguageIds,omitempty" jsonschema:"Translation language IDs to include in the string."`
}

type AddStringInput struct {
	ProjectID   string `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	Name        string `json:"name" jsonschema:"String name/key."`
	Description string `json:"description,omitempty" jsonschema:"Optional string description."`
	Context     string `json:"context,omitempty" jsonschema:"Optional string context."`
	LanguageID  string `json:"languageId,omitempty" jsonschema:"Optional language ID for an initial translation."`
	Translation string `json:"translation,omitempty" jsonschema:"Optional initial translation text. Requires languageId."`
}

type EditStringInput struct {
	ProjectID        string `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	StringID         string `json:"stringId" jsonschema:"String/term ID."`
	Name             string `json:"name,omitempty" jsonschema:"New string name/key. Omit to keep unchanged."`
	Description      string `json:"description,omitempty" jsonschema:"New string description. Omit to keep unchanged unless clearDescription is true."`
	Context          string `json:"context,omitempty" jsonschema:"New string context. Omit to keep unchanged unless clearContext is true."`
	ClearDescription bool   `json:"clearDescription,omitempty" jsonschema:"Set description to null."`
	ClearContext     bool   `json:"clearContext,omitempty" jsonschema:"Set context to null."`
}

type RemoveStringInput struct {
	ProjectID string `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	StringID  string `json:"stringId" jsonschema:"String/term ID to remove."`
}

type SetTranslationInput struct {
	ProjectID  string `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	StringID   string `json:"stringId" jsonschema:"String/term ID."`
	LanguageID string `json:"languageId" jsonschema:"Translation language ID."`
	Text       string `json:"text" jsonschema:"Translation text."`
}

type RemoveTranslationInput struct {
	ProjectID  string `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	StringID   string `json:"stringId" jsonschema:"String/term ID."`
	LanguageID string `json:"languageId" jsonschema:"Translation language ID to remove."`
}

type ExportInput struct {
	ProjectID                string `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	LanguageID               string `json:"languageId" jsonschema:"Language ID to export, for example en or cs."`
	Format                   string `json:"format" jsonschema:"Export format: json, android, apple, or compose."`
	IncludeEmptyTranslations bool   `json:"includeEmptyTranslations,omitempty" jsonschema:"When true, includes terms with empty translations."`
}

type ExportOutput struct {
	ProjectID                string `json:"projectId"`
	LanguageID               string `json:"languageId"`
	Format                   string `json:"format"`
	IncludeEmptyTranslations bool   `json:"includeEmptyTranslations"`
	Content                  string `json:"content"`
}

type ImportInput struct {
	ProjectID        string `json:"projectId,omitempty" jsonschema:"Project ID. Optional when SENTIARY_PROJECT_ID is configured."`
	LanguageID       string `json:"languageId" jsonschema:"Language ID to import into, for example en or cs."`
	Format           string `json:"format" jsonschema:"Import format: json, android, apple, or compose."`
	Content          string `json:"content" jsonschema:"String file content to import."`
	OverrideExisting bool   `json:"overrideExisting,omitempty" jsonschema:"When true, existing translations can be overwritten."`
}

type ProjectBatchInfo struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	TermsLastModified *string  `json:"termsLastModified"`
	Languages         []string `json:"languages"`
}

type ImportResult struct {
	TotalImportedCount int      `json:"totalImportedCount"`
	Skipped            []string `json:"skipped"`
	ImportLanguage     string   `json:"importLanguage"`
}

type Paging struct {
	Items      []Term     `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Size     int     `json:"size"`
	Cursor   *string `json:"cursor"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

type Term struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  *string       `json:"description"`
	Translations []Translation `json:"translations"`
	Updated      string        `json:"updated"`
	AuthorEmail  *string       `json:"authorEmail"`
	AuthorID     *string       `json:"authorId"`
	Statistics   any           `json:"statistics"`
	Context      *string       `json:"context"`
}

type Translation struct {
	LanguageID  string  `json:"languageId"`
	Text        string  `json:"text"`
	Updated     string  `json:"updated"`
	AuthorEmail *string `json:"authorEmail"`
	AuthorID    *string `json:"authorId"`
}

type DeleteOutput struct {
	Deleted bool `json:"deleted"`
}
