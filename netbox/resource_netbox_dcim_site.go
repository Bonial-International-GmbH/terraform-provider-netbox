package netbox

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxDcimSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxDcimSiteCreate,
		Read:   resourceNetboxDcimSiteRead,
		Update: resourceNetboxDcimSiteUpdate,
		Delete: resourceNetboxDcimSiteDelete,
		Exists: resourceNetboxDcimSiteExists,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
			},
		},
	}
}

func resourceNetboxDcimSiteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()

	params := dcim.NewDcimSitesListParams().WithID(&resourceID)

	resources, err := client.Dcim.DcimSitesList(params, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			d.Set("description", resource.Description)
			d.Set("name", *resource.Name)
			d.Set("slug", *resource.Slug)

			return nil
		}
	}

	d.SetId("")

	return nil
}

func resourceNetboxDcimSiteCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	description := d.Get("description").(string)
	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	params := dcim.NewDcimSitesCreateParams()

	params.Data = &models.WritableSite{
		Description: description,
		Slug:        &slug,
		Name:        &name,
	}

	result, err := client.Dcim.DcimSitesCreate(params, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(result.Payload.ID, 10))

	return resourceNetboxDcimSiteRead(d, m)
}

func resourceNetboxDcimSiteUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	description := d.Get("description").(string)
	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	params := dcim.NewDcimSitesPartialUpdateParams().WithID(id)

	params.Data = &models.WritableSite{
		Description: description,
		Slug:        &slug,
		Name:        &name,
	}

	_, err = client.Dcim.DcimSitesPartialUpdate(params, nil)
	if err != nil {
		return err
	}

	return resourceNetboxDcimSiteRead(d, m)
}

func resourceNetboxDcimSiteDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxDcimSiteExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	params := dcim.NewDcimSitesDeleteParams().WithID(id)
	if _, err := client.Dcim.DcimSitesDelete(params, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxDcimSiteExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	id := d.Id()
	params := dcim.NewDcimSitesListParams().WithID(&id)
	sites, err := client.Dcim.DcimSitesList(params, nil)
	if err != nil {
		return resourceExist, err
	}

	for _, site := range sites.Payload.Results {
		if strconv.FormatInt(site.ID, 10) == d.Id() {
			resourceExist = true
		}
	}

	return resourceExist, nil
}
