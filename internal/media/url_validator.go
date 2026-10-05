package media

import (
	"context"
	"net/url"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = urlValidator{}

// urlValidator approximates the API's (Django URLValidator) check: an absolute URL with a
// host and one of the schemes it accepts. The API remains the final authority.
type urlValidator struct{}

func (urlValidator) Description(_ context.Context) string {
	return "value must be an absolute http, https, ftp or ftps URL"
}

func (v urlValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v urlValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	u, err := url.Parse(req.ConfigValue.ValueString())
	if err == nil && u.Host != "" && slices.Contains([]string{"http", "https", "ftp", "ftps"}, u.Scheme) {
		return
	}

	resp.Diagnostics.AddAttributeError(req.Path, "Invalid URL", v.Description(ctx)+", got: "+req.ConfigValue.ValueString())
}
