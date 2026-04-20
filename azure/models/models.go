package models

type RoleDefinitionsResponse struct {
	Value []RoleDefinition `json:"value"`
}

type RoleDefinition struct {
	ID         string                   `json:"id"`
	Name       string                   `json:"name"`
	Properties RoleDefinitionProperties `json:"properties"`
}

type RoleDefinitionProperties struct {
	RoleName         string       `json:"roleName"`
	Description      string       `json:"description"`
	Type             string       `json:"type"`
	Permissions      []Permission `json:"permissions"`
	AssignableScopes []string     `json:"assignableScopes"`
}

type Permission struct {
	Actions     []string `json:"actions"`
	NotActions  []string `json:"notActions"`
	DataActions []string `json:"dataActions"`
}

type RoleAssignmentsResponse struct {
	Value []RoleAssignment `json:"value"`
}

type RoleAssignment struct {
	ID         string                   `json:"id"`
	Properties RoleAssignmentProperties `json:"properties"`
}

type RoleAssignmentProperties struct {
	PrincipalID      string `json:"principalId"`
	RoleDefinitionID string `json:"roleDefinitionId"`
	Scope            string `json:"scope"`
}
