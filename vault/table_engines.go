package vault

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type Engine struct {
	Path string
	Type string
}

func tableEngines() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_engines",
		Description: "Vault secrets engines",
		List: &plugin.ListConfig{
			Hydrate: listEngines,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("path"),
			Hydrate:    getEngine,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path of the secrets engine"},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the secrets engine"},
		},
	}
}

func listEngines(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)

	if err != nil {
		return nil, err
	}

	data, err := conn.Sys().ListMounts()
	for path, _ := range data {
		d.StreamListItem(ctx, &Engine{Type: data[path].Type, Path: path})

	}

	return nil, nil
}

func getEngine(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)

	if err != nil {
		return nil, err
	}

	data, err := conn.Sys().ListMounts()

	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	path := quals["path"].GetStringValue()

	result := data[path]
	if result == nil {
		// TODO: figure out if this is expected to be error
		return nil, nil
	}

	return &Engine{Type: data[path].Type, Path: path}, nil
}
