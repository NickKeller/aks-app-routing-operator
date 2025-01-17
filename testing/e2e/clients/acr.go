package clients

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/Azure/aks-app-routing-operator/testing/e2e/logger"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"golang.org/x/exp/slog"
)

type acr struct {
	id             string
	name           string
	resourceGroup  string
	subscriptionId string
}

func LoadAcr(id arm.ResourceID) *acr {
	return &acr{
		id:             id.String(),
		name:           id.Name,
		resourceGroup:  id.ResourceGroupName,
		subscriptionId: id.SubscriptionID,
	}
}

func NewAcr(ctx context.Context, subscriptionId, resourceGroup, name, location string) (*acr, error) {
	name = nonAlphanumericRegex.ReplaceAllString(name, "")

	lgr := logger.FromContext(ctx).With("name", name, "resourceGroup", resourceGroup, "location", location)
	ctx = logger.WithContext(ctx, lgr)
	lgr.Info("starting to create acr")
	defer lgr.Info("finished creating acr")

	cred, err := getAzCred()
	if err != nil {
		return nil, fmt.Errorf("getting az credentials: %w", err)
	}

	factory, err := armcontainerregistry.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("creating client factory: %w", err)
	}

	r := &armcontainerregistry.Registry{
		Location: to.Ptr(location),
		SKU: &armcontainerregistry.SKU{
			Name: to.Ptr(armcontainerregistry.SKUNameBasic),
		},
	}
	poller, err := factory.NewRegistriesClient().BeginCreate(ctx, resourceGroup, name, *r, nil)
	if err != nil {
		return nil, fmt.Errorf("starting to create registry: %w", err)
	}

	result, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("waiting for registry creation to complete: %w", err)
	}

	return &acr{
		id:             *result.ID,
		name:           *result.Name,
		resourceGroup:  resourceGroup,
		subscriptionId: subscriptionId,
	}, nil
}

func (a *acr) GetName() string {
	return a.name
}

func (a *acr) GetId() string {
	return a.id
}

func (a *acr) BuildAndPush(ctx context.Context, imageName, dockerfilePath string) error {
	lgr := logger.FromContext(ctx).With("image", imageName, "name", a.name, "resourceGroup", a.resourceGroup, "subscriptionId", a.subscriptionId)
	ctx = logger.WithContext(ctx, lgr)
	lgr.Info("starting to build and push image")
	defer lgr.Info("finished building and pushing image")

	// Ideally, we'd use the sdk to build and push the image but I couldn't get it working.
	// I matched everything on the az cli but wasn't able to get it working with the sdk.
	// https://github.com/Azure/azure-cli/blob/5f9a8fa25cc1c980ebe5e034bd419c95a1c578e2/src/azure-cli/azure/cli/command_modules/acr/build.py#L25
	cmd := exec.Command("az", "acr", "build", "--registry", a.name, "--image", imageName, dockerfilePath)
	cmd.Stdout = newLogWriter(lgr, "building and pushing acr image: ", nil)
	cmd.Stderr = newLogWriter(lgr, "building and pushing acr image: ", to.Ptr(slog.LevelError))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("starting build and push command: %w", err)
	}

	return nil
}
