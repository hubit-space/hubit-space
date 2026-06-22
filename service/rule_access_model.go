package middleware

import (
	"hubit-space/service/handler"
	"hubit-space/service/repository"
	"hubit-space/service/utility"
	"log"

	"github.com/gin-gonic/gin"
)

func AccessDataFromCreatedByCode(repo repository.RuleAccessRepository, types string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCode, _ := utility.GetHeader(c, "RoleCode")
		userCode, _ := utility.GetHeader(c, "UserCode")

		// fetching rule access details
		ruleAccess, err := handler.NewRuleAccessHandler(repo).GetRuleAccessDetail(roleCode, types)
		if err != nil {
			log.Println("Error fetching rule access:", err)
			return
		}

		if ruleAccess["rule_code"] == "OWN" {
			query := c.Request.URL.Query()

			existingFilter := query.Get("filter")
			if existingFilter == "" {
				query.Set("filter", "created_by_code:"+userCode)
			} else {
				query.Set("filter", existingFilter+",created_by_code:"+userCode)
			}

			c.Request.URL.RawQuery = query.Encode()
		}

		c.Next()
	}
}
