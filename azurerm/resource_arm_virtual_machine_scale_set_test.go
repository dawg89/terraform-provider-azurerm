package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMVirtualMachineScaleSet_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basic(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),

					// single placement group should default to true
					testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup("azurerm_virtual_machine_scale_set.test", true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicPublicIP(t *testing.T) {
	resource_name := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicPublicIP(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resource_name),
					testCheckAzureRMVirtualMachineScaleSetIsPrimary(resource_name, true),
					testCheckAzureRMVirtualMachineScaleSetPublicIPName(resource_name, "TestPublicIPConfiguration"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicAcceleratedNetworking(t *testing.T) {
	resource_name := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicAcceleratedNetworking(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resource_name),
					testCheckAzureRMVirtualMachineScaleSetAcceleratedNetworking(resource_name, true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_bootDiagnostic(t *testing.T) {
	resource_name := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_bootDiagnostic(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "boot_diagnostics.0.enabled", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_networkSecurityGroup(t *testing.T) {
	resource_name := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_networkSecurityGroup(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resource_name),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicWindows(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicWindows(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),

					// single placement group should default to true
					testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup("azurerm_virtual_machine_scale_set.test", true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_singlePlacementGroupFalse(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_singlePlacementGroupFalse(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
					testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup("azurerm_virtual_machine_scale_set.test", false),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_linuxUpdated(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineScaleSet_linux(ri, location)
	updatedConfig := testAccAzureRMVirtualMachineScaleSet_linuxUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_customDataUpdated(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineScaleSet_linux(ri, location)
	updatedConfig := testAccAzureRMVirtualMachineScaleSet_linuxCustomDataUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(ri, testLocation(), "Standard_D1_v2")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk_resize(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(ri, testLocation(), "Standard_D1_v2")
	postConfig := testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(ri, testLocation(), "Standard_D2_v2")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicLinux_managedDiskNoName(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDiskNoName(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_basicLinux_disappears(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_basic(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
					testCheckAzureRMVirtualMachineScaleSetDisappears("azurerm_virtual_machine_scale_set.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_planManagedDisk(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_planManagedDisk(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_customImage(t *testing.T) {
	ri := acctest.RandInt()
	resourceGroup := fmt.Sprintf("acctestRG-%d", ri)
	userName := fmt.Sprintf("testadmin%d", ri)
	password := fmt.Sprintf("Password1234!%d", ri)
	hostName := fmt.Sprintf("tftestcustomimagesrc%d", ri)
	sshPort := "22"
	config := testAccAzureRMVirtualMachineScaleSet_customImage(ri, testLocation(), userName, password, hostName)
	preConfig := testAccAzureRMImage_standaloneImage_setup(ri, userName, password, hostName, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				//need to create a vm and then reference it in the image creation
				Config:  preConfig,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureVMExists("azurerm_virtual_machine.testsource", true),
					testGeneralizeVMImage(resourceGroup, "testsource", userName, password, hostName, sshPort, testLocation()),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
					testCheckAzureRMImageExists("azurerm_image.test", true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_loadBalancer(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
					testCheckAzureRMVirtualMachineScaleSetHasLoadbalancer("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_loadBalancerManagedDataDisks(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplateManagedDataDisks(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
					testCheckAzureRMVirtualMachineScaleSetHasDataDisks("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_overprovision(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetOverProvisionTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
					testCheckAzureRMVirtualMachineScaleSetOverprovision("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_MSI(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetMSITemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttrSet(resourceName, "identity.0.principal_id"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_extension(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetExtensionTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
					testCheckAzureRMVirtualMachineScaleSetExtension("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_extensionUpdate(t *testing.T) {
	resourceName := "azurerm_virtual_machine_scale_set.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMVirtualMachineScaleSetExtensionTemplate(ri, location)
	updatedConfig := testAccAzureRMVirtualMachineScaleSetExtensionTemplateUpdated(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(resourceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists(resourceName),
					testCheckAzureRMVirtualMachineScaleSetExtension(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_multipleExtensions(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetMultipleExtensionsTemplate(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
					testCheckAzureRMVirtualMachineScaleSetExtension("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_osDiskTypeConflict(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_osDiskTypeConflict(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Conflict between `vhd_containers`"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_NonStandardCasing(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSetNonStandardCasing(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
				),
			},
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMVirtualMachineScaleSet_multipleNetworkProfiles(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachineScaleSet_multipleNetworkProfiles(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineScaleSetExists("azurerm_virtual_machine_scale_set.test"),
				),
			},
		},
	})
}

func testGetAzureRMVirtualMachineScaleSet(s *terraform.State, resourceName string) (result *compute.VirtualMachineScaleSet, err error) {
	// Ensure we have enough information in state to look up in API
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", resourceName)
	}

	// Name of the actual scale set
	name := rs.Primary.Attributes["name"]

	resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
	if !hasResourceGroup {
		return nil, fmt.Errorf("Bad: no resource group found in state for virtual machine: scale set %s", name)
	}

	client := testAccProvider.Meta().(*ArmClient).vmScaleSetClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	vmss, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("Bad: Get on vmScaleSetClient: %+v", err)
	}

	if vmss.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Bad: VirtualMachineScaleSet %q (resource group: %q) does not exist", name, resourceGroup)
	}

	return &vmss, err
}

func testCheckAzureRMVirtualMachineScaleSetExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		return err
	}
}

func testCheckAzureRMVirtualMachineScaleSetDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: scale set %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).vmScaleSetClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Delete on vmScaleSetClient: %+v", err)
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Bad: Delete on vmScaleSetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).vmScaleSetClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_machine_scale_set" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual Machine Scale Set still exists:\n%#v", resp.VirtualMachineScaleSetProperties)
		}
	}

	return nil
}

func testCheckAzureRMVirtualMachineScaleSetHasLoadbalancer(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		ip := (*n)[0].IPConfigurations
		if ip == nil || len(*ip) == 0 {
			return fmt.Errorf("Bad: Could not get ip configurations for scale set %v", name)
		}

		pools := (*ip)[0].LoadBalancerBackendAddressPools
		if pools == nil || len(*pools) == 0 {
			return fmt.Errorf("Bad: Load balancer backend pools is empty for scale set %v", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetIsPrimary(name string, boolean bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		ip := (*n)[0].IPConfigurations
		if ip == nil || len(*ip) == 0 {
			return fmt.Errorf("Bad: Could not get ip configurations for scale set %v", name)
		}

		primary := *(*ip)[0].Primary
		if primary != boolean {
			return fmt.Errorf("Bad: Primary set incorrectly for scale set %v\n Wanted: %+v Received: %+v", name, boolean, primary)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetPublicIPName(name, publicIPName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		ip := (*n)[0].IPConfigurations
		if ip == nil || len(*ip) == 0 {
			return fmt.Errorf("Bad: Could not get ip configurations for scale set %v", name)
		}

		publicIPConfigName := *(*ip)[0].PublicIPAddressConfiguration.Name
		if publicIPConfigName != publicIPName {
			return fmt.Errorf("Bad: Public IP Config Name set incorrectly for scale set %v\n Wanted: %+v Received: %+v", name, publicIPName, publicIPConfigName)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetAcceleratedNetworking(name string, boolean bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get network interface configurations for scale set %v", name)
		}

		acceleratedNetworking := *(*n)[0].EnableAcceleratedNetworking
		if acceleratedNetworking != boolean {
			return fmt.Errorf("Bad: Primary set incorrectly for scale set %v\n Wanted: %+v Received: %+v", name, boolean, acceleratedNetworking)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetOverprovision(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		if *resp.Overprovision {
			return fmt.Errorf("Bad: Overprovision should have been false for scale set %v", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetSinglePlacementGroup(name string, expectedSinglePlacementGroup bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		if *resp.SinglePlacementGroup != expectedSinglePlacementGroup {
			return fmt.Errorf("Bad: Overprovision should have been %t for scale set %v", expectedSinglePlacementGroup, name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetMSI(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		identityType := resp.Identity.Type
		if identityType != "systemAssigned" {
			return fmt.Errorf("Bad: Identity Type is not systemAssigned for scale set %v", name)
		}

		principalID := *resp.Identity.PrincipalID
		if len(principalID) == 0 {
			return fmt.Errorf("Bad: Could not get principal_id for scale set %v", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetExtension(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resp, err := testGetAzureRMVirtualMachineScaleSet(s, name)
		if err != nil {
			return err
		}

		n := resp.VirtualMachineProfile.ExtensionProfile.Extensions
		if n == nil || len(*n) == 0 {
			return fmt.Errorf("Bad: Could not get extensions for scale set %v", name)
		}

		return nil
	}
}

func testCheckAzureRMVirtualMachineScaleSetHasDataDisks(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual machine: scale set %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).vmScaleSetClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on vmScaleSetClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: VirtualMachineScaleSet %q (resource group: %q) does not exist", name, resourceGroup)
		}

		storageProfile := resp.VirtualMachineProfile.StorageProfile.DataDisks
		if storageProfile == nil || len(*storageProfile) == 0 {
			return fmt.Errorf("Bad: Could not get data disks configurations for scale set %v", name)
		}

		return nil
	}
}

func testAccAzureRMVirtualMachineScaleSet_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }
  
  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_basicPublicIP(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
				primary = true
				public_ip_address_configuration {
					name = "TestPublicIPConfiguration"
					domain_name_label = "test-domain-label-%d"
					idle_timeout = 4
				}
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_basicAcceleratedNetworking(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D4_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%d"
      primary = true
			accelerated_networking = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_bootDiagnostic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

	boot_diagnostics {
			storage_uri = "${azurerm_storage_account.test.primary_blob_endpoint}"
	}

  network_profile {
      name = "TestNetworkProfile-%d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_networkSecurityGroup(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%d"
      primary = true
			network_security_group_id = "${azurerm_network_security_group.test.id}"
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
				primary = true
				public_ip_address_configuration {
					name = "TestPublicIPConfiguration"
					domain_name_label = "test-domain-label-%d"
					idle_timeout = 4
				}
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_basicWindows(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
	name                     = "accsa%d"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  os_profile_windows_config {
    enable_automatic_upgrades = false
    provision_vm_agent = true

    winrm {
      protocol = "http"
    }
  }

  network_profile {
      name = "TestNetworkProfile-%d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_singlePlacementGroupFalse(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%[1]d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%[1]d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%[1]d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%[1]d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  single_placement_group = false

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%[1]d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%[1]d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = ""
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location)
}

func testAccAzureRMVirtualMachineScaleSet_linux(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.0.0.0/8"]
}
resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}
resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}
resource "azurerm_public_ip" "test" {
  name                         = "acctestpip-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  public_ip_address_allocation = "static"
}
resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}
resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  upgrade_policy_mode = "Automatic"
  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = "1"
  }
  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
    admin_password       = "password"
    custom_data          = "custom data!"
  }
  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/ubuntu/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }
  network_profile {
    name    = "TestNetworkProfile"
    primary = true
    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
    }
  }
  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    os_type        = "linux"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }
  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_linuxUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.0.0.0/8"]
}
resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}
resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}
resource "azurerm_public_ip" "test" {
  name                         = "acctestpip-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  public_ip_address_allocation = "static"
}
resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}
resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  upgrade_policy_mode = "Automatic"
  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = "1"
  }
  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
    admin_password       = "password"
    custom_data          = "custom data!"
  }
  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/ubuntu/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }
  network_profile {
    name    = "TestNetworkProfile"
    primary = true
    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
    }
  }
  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    os_type        = "linux"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }
  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
  tags {
    ThisIs = "a test"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_linuxCustomDataUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.0.0.0/8"]
}
resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}
resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}
resource "azurerm_public_ip" "test" {
  name                         = "acctestpip-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  public_ip_address_allocation = "static"
}
resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  frontend_ip_configuration {
    name                 = "ip-address"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
resource "azurerm_lb_backend_address_pool" "test" {
  name                = "acctestbap-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}
resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctestvmss-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  upgrade_policy_mode = "Automatic"
  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = "1"
  }
  os_profile {
    computer_name_prefix = "prefix"
    admin_username       = "ubuntu"
    admin_password       = "password"
    custom_data          = "updated custom data!"
  }
  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/ubuntu/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }
  network_profile {
    name    = "TestNetworkProfile"
    primary = true
    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
    }
  }
  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    os_type        = "linux"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }
  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
    name = "TestNetworkProfile-%d"
    primary = true
    ip_configuration {
      name = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name 		  = ""
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDisk_withZones(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
    name = "acctvmss-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    upgrade_policy_mode = "Manual"
    zones = ["1", "2"]

    sku {
        name = "Standard_D1_v2"
        tier = "Standard"
        capacity = 2
    }

    os_profile {
        computer_name_prefix = "testvm-%d"
        admin_username = "myadmin"
        admin_password = "Passwword1234"
    }

    network_profile {
        name = "TestNetworkProfile-%d"
        primary = true
        ip_configuration {
            name = "TestIPConfiguration"
            subnet_id = "${azurerm_subnet.test.id}"
        }
    }

    storage_profile_os_disk {
        caching = "ReadWrite"
        create_option = "FromImage"
        managed_disk_type = "Standard_LRS"
    }

    storage_profile_image_reference {
        publisher = "Canonical"
        offer = "UbuntuServer"
        sku = "16.04-LTS"
        version = "latest"
    }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_basicWindows_managedDisk(rInt int, location string, vmSize string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "%s"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  os_profile_windows_config {
	enable_automatic_upgrades = false
	provision_vm_agent = true

	additional_unattend_config {
	  pass = "oobeSystem"
	  component = "Microsoft-Windows-Shell-Setup"
	  setting_name = "AutoLogon"
	  content = "<AutoLogon><Username>myadmin</Username><Password><Value>Passwword1234</Value></Password><Enabled>true</Enabled><LogonCount>1</LogonCount></AutoLogon>"
	}

	additional_unattend_config {
      pass = "oobeSystem"
      component = "Microsoft-Windows-Shell-Setup"
      setting_name = "FirstLogonCommands"
      content = "<FirstLogonCommands><SynchronousCommand><CommandLine>shutdown /r /t 0 /c \"initial reboot\"</CommandLine><Description>reboot</Description><Order>1</Order></SynchronousCommand></FirstLogonCommands>"
    }
  }

  network_profile {
    name = "TestNetworkProfile-%d"
    primary = true
    ip_configuration {
      name = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name 		  = ""
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, vmSize, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_basicLinux_managedDiskNoName(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
    name = "TestNetworkProfile-%d"
    primary = true
    ip_configuration {
      name = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_lb_nat_pool" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  name                           = "ssh"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  protocol                       = "Tcp"
  frontend_port_start            = 50000
  frontend_port_end              = 50119
  backend_port                   = 22
  frontend_ip_configuration_name = "default"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
      load_balancer_inbound_nat_rules_ids    = ["${azurerm_lb_nat_pool.test.id}"]
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetOverProvisionTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetMSITemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  identity {
    type     = "systemAssigned"
  }
  
  extension {
    name                       = "MSILinuxExtension"
    publisher                  = "Microsoft.ManagedIdentity"
    type                       = "ManagedIdentityExtensionForLinux"
    type_handler_version       = "1.0"
    settings                   = "{\"port\": 50342}"
  }  

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetExtensionTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/myadmin/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = <<SETTINGS
		{
			"commandToExecute": "echo $HOSTNAME"
		}
SETTINGS

    protected_settings = <<SETTINGS
		{
			"storageAccountName": "${azurerm_storage_account.test.name}",
			"storageAccountKey": "${azurerm_storage_account.test.primary_access_key}"
		}
SETTINGS
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetExtensionTemplateUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/myadmin/.ssh/authorized_keys"
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDCsTcryUl51Q2VSEHqDRNmceUFo55ZtcIwxl2QITbN1RREti5ml/VTytC0yeBOvnZA4x4CFpdw/lCDPk0yrH9Ei5vVkXmOrExdTlT3qI7YaAzj1tUVlBd4S6LX1F7y6VLActvdHuDDuXZXzCDd/97420jrDfWZqJMlUK/EmCE5ParCeHIRIvmBxcEnGfFIsw8xQZl0HphxWOtJil8qsUWSdMyCiJYYQpMoMliO99X40AUc4/AlsyPyT5ddbKk08YrZ+rKDVHF7o29rh4vi5MmHkVgVQHKiKybWlHq+b71gIAUQk9wrJxD+dqt4igrmDSpIjfjwnd+l5UIn5fJSO5DYV4YT/4hwK7OKmuo7OFHD0WyY5YnkYEMtFgzemnRBdE8ulcT60DQpVgRMXFWHvhyCWy0L6sgj1QWDZlLpvsIvNfHsyhKFMG1frLnMt/nP0+YCcfg+v1JYeCKjeoJxB8DWcRBsjzItY0CGmzP8UYZiYKl/2u+2TgFS5r7NWH11bxoUzjKdaa1NLw+ieA8GlBFfCbfWe6YVB9ggUte4VtYFMZGxOjS2bAiYtfgTKFJv+XqORAwExG6+G2eDxIDyo80/OA9IG7Xv/jwQr7D6KDjDuULFcN/iTxuttoKrHeYz1hf5ZQlBdllwJHYx6fK2g8kha6r2JIQKocvsAXiiONqSfw== hello@world.com"
    }
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = <<SETTINGS
		{
			"commandToExecute": "echo $HOSTNAME",
			"timestamp": 12345679955
		}
SETTINGS

    protected_settings = <<SETTINGS
		{
			"storageAccountName": "${azurerm_storage_account.test.name}",
			"storageAccountKey": "${azurerm_storage_account.test.primary_access_key}"
		}
SETTINGS
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetMultipleExtensionsTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "os-disk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  extension {
    name                       = "CustomScript"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "CustomScript"
    type_handler_version       = "2.0"
    auto_upgrade_minor_version = true

    settings = <<SETTINGS
		{
			"commandToExecute": "echo $HOSTNAME"
		}
SETTINGS

    protected_settings = <<SETTINGS
		{
			"storageAccountName": "${azurerm_storage_account.test.name}",
			"storageAccountKey": "${azurerm_storage_account.test.primary_access_key}"
		}
SETTINGS
  }

  extension {
    name                       = "Docker"
    publisher                  = "Microsoft.Azure.Extensions"
    type                       = "DockerExtension"
    type_handler_version       = "1.0"
    auto_upgrade_minor_version = true
  }
}

`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_osDiskTypeConflict(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_D1_v2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile-%d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
    vhd_containers    = ["should_cause_conflict"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetLoadBalancerTemplateManagedDataDisks(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "test"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_A0"
    tier     = "Standard"
    capacity = 1
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name                                   = "TestIPConfiguration"
      subnet_id                              = "${azurerm_subnet.test.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.test.id}"]
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_data_disk {
    lun               = 0
    caching           = "ReadWrite"
    create_option     = "Empty"
    disk_size_gb      = 10
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}

`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSetNonStandardCasing(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_A0"
    tier     = "standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile-%d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }

  storage_profile_os_disk {
    name           = "osDiskProfile"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_planManagedDisk(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  plan {
    name = "os"
    product = "rancheros"
    publisher = "rancher"
  }

  storage_profile_os_disk {
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "rancher"
    offer     = "rancheros"
    sku       = "os"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_customImage(rInt int, location string, userName string, password string, hostName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctpip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
  domain_name_label            = "%s"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "blob"
}

resource "azurerm_virtual_machine" "testsource" {
  name                  = "testsource"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.testsource.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "mdimagetestsource"
    admin_username = "%s"
    admin_password = "%s"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_image" "test" {
  name                = "accteste-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    size_gb  = 30
    caching  = "None"
  }

  tags {
    environment = "Dev"
    cost-center = "Ops"
  }
}


 data "azurerm_image" "test" {
  name = "${azurerm_image.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

	storage_profile_image_reference {
   id = "${data.azurerm_image.test.id}"
	}

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "TestNetworkProfile-%d"
      primary = true
      ip_configuration {
        name = "TestIPConfiguration"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    caching       = "ReadWrite"
    create_option = "FromImage"
    managed_disk_type = "Standard_LRS"
  }
}
`, rInt, location, rInt, rInt, rInt, hostName, rInt, rInt, userName, password, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachineScaleSet_multipleNetworkProfiles(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_storage_account" "test" {
    name = "accsa%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    account_tier = "Standard"
    account_replication_type = "LRS"

    tags {
        environment = "staging"
    }
}

resource "azurerm_storage_container" "test" {
    name = "vhds"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_name = "${azurerm_storage_account.test.name}"
    container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name = "acctvmss-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name = "Standard_D1_v2"
    tier = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username = "myadmin"
    admin_password = "Passwword1234"
  }

  network_profile {
      name = "primary-%d"
      primary = true
      ip_configuration {
        name = "primary"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  network_profile {
      name = "secondary-%d"
      primary = false
      ip_configuration {
        name = "secondary"
        subnet_id = "${azurerm_subnet.test.id}"
      }
  }

  storage_profile_os_disk {
    name = "osDiskProfile"
    caching       = "ReadWrite"
    create_option = "FromImage"
    vhd_containers = ["${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"]
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}
