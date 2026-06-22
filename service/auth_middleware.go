package handler

import (
	"hubit-space/service/repository"
	"log"
)

type RuleAccessHandler struct {
	repo repository.RuleAccessRepository
}

func NewRuleAccessHandler(repo repository.RuleAccessRepository) *RuleAccessHandler {
	return &RuleAccessHandler{
		repo: repo,
	}
}

func (h *RuleAccessHandler) GetRuleAccessDetail(roleCode string, menuCode string) (map[string]any, error) {
	ruleAccess, err := h.repo.GetRuleAccessDetail(roleCode, menuCode)
	if err != nil {
		log.Println("Error fetching rule access:", err)
		return nil, err
	}

	response := map[string]any{
		"role_code": ruleAccess.RoleCode,
		"rule_code": ruleAccess.RuleCode,
	}

	return response, nil
}
