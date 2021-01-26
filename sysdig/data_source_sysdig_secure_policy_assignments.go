package sysdig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecurePolicyAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigSecurePolicyAssignmentsRead,
		Schema: map[string]*schema.Schema{
			"policy_assignments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"registry": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_ids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"whitelist_ids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSysdigSecurePolicyAssignmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diagz diag.Diagnostics

	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := "default"

	d.SetId(name)

	providerBundle, err := client.GetPolicyAssignments(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	policyAssignments := make([]interface{}, len(providerBundle.Items), len(providerBundle.Items))

	for i, item := range providerBundle.Items {
		bundleItem := make(map[string]interface{})

		bundleItem["registry"] = item.Registry
		bundleItem["id"] = item.ID
		bundleItem["repository"] = item.Repository

		// policy ids
		policyIds := []string{}
		for _, policy := range item.Policies {
			policyIds = append(policyIds, policy)
		}
		bundleItem["policy_ids"] = policyIds

		// whitelist ids
		whitelistIds := []string{}
		for _, item := range item.Whitelist {
			whitelistIds = append(whitelistIds, item)
		}
		bundleItem["whitelist_ids"] = whitelistIds

		bundleItem["tag"] = item.Image.Value

		policyAssignments[i] = bundleItem
	}

	if err = d.Set("policy_assignments", policyAssignments); err != nil {
		return diag.FromErr(err)
	}

	return diagz
}