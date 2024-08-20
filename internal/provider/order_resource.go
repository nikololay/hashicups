package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &orderResource{}
	_ resource.ResourceWithConfigure = &orderResource{}
)

func NewOrderResource() resource.Resource {
	return &orderResource{}
}

type orderResource struct {
	client *hashicups.Client
}

type orderResourceModel struct {
	ID          types.String     `tfsdk:"id"`
	Items       []orderItemModel `tfsdk:"items"`
	LastUpdated types.String     `tfsdk:"last_updated"`
}

type orderItemModel struct {
	Coffee   orderItemCoffeeModel `tfsdk:"coffee"`
	Quantity types.Int64          `tfsdk:"quantity"`
}

type orderItemCoffeeModel struct {
	ID          types.Int64   `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Teaser      types.String  `tfsdk:"teaser"`
	Description types.String  `tfsdk:"description"`
	Price       types.Float64 `tfsdk:"price"`
	Image       types.String  `tfsdk:"image"`
}

// Configure implements resource.ResourceWithConfigure.
func (o *orderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hashicups.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.",
		)

		return
	}

	o.client = client
}

// Metadata implements resource.Resource.
func (o *orderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_order"
}

// Schema defines the schema for the resource.
func (r *orderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"items": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"quantity": schema.Int64Attribute{
							Required: true,
						},
						"coffee": schema.SingleNestedAttribute{
							Required: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Required: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
								"teaser": schema.StringAttribute{
									Computed: true,
								},
								"description": schema.StringAttribute{
									Computed: true,
								},
								"price": schema.Float64Attribute{
									Computed: true,
								},
								"image": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (o *orderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan orderResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate API request body from plan
	var items []hashicups.OrderItem
	for _, item := range plan.Items {
		items = append(items, hashicups.OrderItem{
			Coffee: hashicups.Coffee{
				ID: int(item.Coffee.ID.ValueInt64()),
			},
			Quantity: int(item.Quantity.ValueInt64()),
		})
	}

	// create new order
	order, err := o.client.CreateOrder(items)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create order", err.Error())
		return
	}

	// map response body to schema and populated Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(order.ID))
	for orderItemIndex, orderItem := range order.Items {
		plan.Items[orderItemIndex] = orderItemModel{
			Coffee: orderItemCoffeeModel{
				ID:          types.Int64Value(int64(orderItem.Coffee.ID)),
				Name:        types.StringValue(orderItem.Coffee.Name),
				Teaser:      types.StringValue(orderItem.Coffee.Teaser),
				Description: types.StringValue(orderItem.Coffee.Description),
				Price:       types.Float64Value(orderItem.Coffee.Price),
				Image:       types.StringValue(orderItem.Coffee.Image),
			},
			Quantity: types.Int64Value(int64(orderItem.Quantity)),
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read implements resource.Resource.
func (o *orderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// get the current state
	var state orderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get refreshed state from the API
	order, err := o.client.GetOrder(state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError("Unable to fetch order", err.Error())
		return
	}

	// overwrite items with the refreshed state
	state.Items = []orderItemModel{}
	for _, orderItem := range order.Items {
		state.Items = append(state.Items, orderItemModel{
			Coffee: orderItemCoffeeModel{
				ID:          types.Int64Value(int64(orderItem.Coffee.ID)),
				Name:        types.StringValue(orderItem.Coffee.Name),
				Teaser:      types.StringValue(orderItem.Coffee.Teaser),
				Description: types.StringValue(orderItem.Coffee.Description),
				Price:       types.Float64Value(orderItem.Coffee.Price),
				Image:       types.StringValue(orderItem.Coffee.Image),
			},
			Quantity: types.Int64Value(int64(orderItem.Quantity)),
		})
	}

	// set the refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update implements resource.Resource.
func (o *orderResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}

// Delete implements resource.Resource.
func (o *orderResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("unimplemented")
}
