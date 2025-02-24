package oktaasa

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccProject(t *testing.T) {
	project := &Project{}
	projectName := "test-acc-project"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccProjectCheckDestroy(project),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("oktaasa_project.test", project),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "project_name", projectName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_uid", "60120",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_gid", "63020",
					),
				),
			},
			{
				Config: testAccProjectUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("oktaasa_project.test", project),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "project_name", projectName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_uid", "61200",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_gid", "63400",
					),
				),
			},
			{
				Config: testAccProjectCreateConfig2,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("oktaasa_project.test", project),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "project_name", projectName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_uid", "0",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_gid", "0",
					),
				),
			},
			{
				Config: testAccProjectCreateConfigFull,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("oktaasa_project.test", project),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "project_name", projectName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_uid", "0",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_gid", "0",
					),
				),
			},
		},
	})
}

func testAccProjectCheckExists(rn string, p *Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		// resource ID is project name
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/projects/"+rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		err = json.Unmarshal(resp.Body(), p)
		if err != nil {
			return fmt.Errorf("error unmarshaling data source response: %s", err)
		}

		return nil
	}
}

func testAccProjectCheckDestroy(p *Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/projects/"+p.Name)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		status := resp.StatusCode()
		deleted, err := checkSoftDelete(resp.Body())
		if err != nil {
			return fmt.Errorf("error while checking deleted status: %s", err)
		}

		if status == 200 && !deleted {
			return fmt.Errorf("project still exists")
		}

		return nil
	}
}

const testAccProjectCreateConfig = `
resource "oktaasa_project" "test" {
    project_name = "test-acc-project"
  	next_unix_uid = 60120
  	next_unix_gid = 63020
}`

const testAccProjectCreateConfig2 = `
resource "oktaasa_project" "test" {
    project_name = "test-acc-project"
}`

const testAccProjectCreateConfigFull = `
resource "oktaasa_project" "test" {
    project_name              = "test-acc-project"
  	next_unix_uid             = 60120
  	next_unix_gid             = 63020
    create_server_users       = true
    force_shared_ssh_users    = true
    forward_traffic           = true
    rdp_session_recording     = true
    require_preauthorization  = true
    shared_admin_user_name    = "sauser"
    shared_standard_user_name = "ssuser"
    ssh_session_recording     = true
}`

const testAccProjectUpdateConfig = `
resource "oktaasa_project" "test" {
    project_name = "test-acc-project"
  	next_unix_uid = 61200
  	next_unix_gid = 63400
}`
