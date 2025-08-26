package mdbcommon

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/yandex-cloud/terraform-provider-yandex/pkg/datasize"
	"github.com/yandex-cloud/terraform-provider-yandex/pkg/validate"
	utils "github.com/yandex-cloud/terraform-provider-yandex/pkg/wrappers"
	"github.com/yandex-cloud/terraform-provider-yandex/yandex-framework/provider/config"
	"google.golang.org/genproto/googleapis/type/timeofday"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ExpandBackupWindow(ctx context.Context, bw types.Object, diags *diag.Diagnostics) *timeofday.TimeOfDay {
	backupWindow := &BackupWindow{}
	diags.Append(bw.As(ctx, backupWindow, datasize.UnhandledOpts)...)
	if diags.HasError() {
		return nil
	}
	rs := &timeofday.TimeOfDay{
		Hours:   int32(backupWindow.Hours.ValueInt64()),
		Minutes: int32(backupWindow.Minutes.ValueInt64()),
	}
	return rs
}

func ExpandResources[V any, T resourceModel[V]](ctx context.Context, o types.Object, diags *diag.Diagnostics) T {
	if !utils.IsPresent(o) {
		return nil
	}

	d := &Resource{}
	diags.Append(o.As(ctx, d, baseOptions)...)
	if diags.HasError() {
		return nil
	}

	rs := T(new(V))
	rs.SetResourcePresetId(d.ResourcePresetId.ValueString())
	rs.SetDiskSize(datasize.ToBytes(d.DiskSize.ValueInt64()))
	if utils.IsPresent(d.DiskTypeId) {
		rs.SetDiskTypeId(d.DiskTypeId.ValueString())
	}

	return rs
}

var environments = map[string]int32{
	"ENVIRONMENT_UNSPECIFIED": 0,
	"PRODUCTION":              1,
	"PRESTABLE":               2,
}

func ExpandEnvironment[T ~int32](_ context.Context, e types.String, diags *diag.Diagnostics) T {
	if !utils.IsPresent(e) {
		return 0
	}

	v, ok := environments[e.ValueString()]
	if !ok || v == 0 {
		allowedEnvs := make([]string, 0, len(environments))
		for k, v := range environments {
			if v == 0 {
				continue
			}
			allowedEnvs = append(allowedEnvs, k)
		}

		diags.AddError(
			"Failed to parse environment",
			fmt.Sprintf("Error while parsing value for 'environment'. Value must be one of `%s`, not `%s`", strings.Join(allowedEnvs, "`, `"), e),
		)

		return 0
	}
	return T(v)
}

func ExpandLabels(ctx context.Context, labels types.Map, diags *diag.Diagnostics) map[string]string {
	var lMap map[string]string
	if utils.IsPresent(labels) {
		diags.Append(labels.ElementsAs(ctx, &lMap, false)...)
		if diags.HasError() {
			return nil
		}
	}
	return lMap
}

func ExpandClusterMaintenanceWindow[
	V any,
	VW any,
	VA any,

	WD ~int32,

	T maintenanceWindow[
		V,
		VW, VA,
		WD,
		W, A,
	],

	W weeklyMaintenanceWindow[VW, WD],
	A anytimeMaintenanceWindow[VA],
](ctx context.Context, mw types.Object, diags *diag.Diagnostics) T {
	if !utils.IsPresent(mw) {
		return *new(T)
	}

	out := T(new(V))
	var mwConf MaintenanceWindow

	diags.Append(mw.As(ctx, &mwConf, datasize.DefaultOpts)...)
	if diags.HasError() {
		return *new(T)
	}

	if mwType := mwConf.Type.ValueString(); mwType == anytimeType {
		out.SetAnytime(new(VA))
	} else if mwType == weeklyType {
		mwDay, mwHour := mwConf.Day.ValueString(), mwConf.Hour.ValueInt64()
		day := weekdayNames[mwDay]

		w := W(new(VW))
		w.SetDay(WD(day))
		w.SetHour(mwHour)
		out.SetWeeklyMaintenanceWindow(w)
	} else {
		diags.AddError(
			"Failed to expand maintenance window.",
			fmt.Sprintf("maintenance_window.type should be %s or %s", anytimeType, weeklyType),
		)
		return *new(T)
	}

	return out
}

func ExpandBoolWrapper(_ context.Context, b types.Bool, _ *diag.Diagnostics) *wrapperspb.BoolValue {
	if b.IsNull() || b.IsUnknown() {
		return nil
	}

	return wrapperspb.Bool(b.ValueBool())
}

func ExpandStringWrapper(_ context.Context, s types.String, _ *diag.Diagnostics) *wrapperspb.StringValue {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}

	return wrapperspb.String(s.ValueString())
}

func ExpandSecurityGroupIds(ctx context.Context, sg types.Set, diags *diag.Diagnostics) []string {
	var securityGroupIds []string
	if !(sg.IsUnknown() || sg.IsNull()) {
		securityGroupIds = make([]string, len(sg.Elements()))
		diags.Append(sg.ElementsAs(ctx, &securityGroupIds, false)...)
		if diags.HasError() {
			return nil
		}
	}

	return securityGroupIds
}

func ExpandFolderId(ctx context.Context, f types.String, providerConfig *config.State, diags *diag.Diagnostics) string {
	folderID, d := validate.FolderID(f, providerConfig)
	diags.Append(d)
	return folderID
}

func ExpandInt64Wrapper(ctx context.Context, in types.Int64, diags *diag.Diagnostics) *wrapperspb.Int64Value {
	var w *wrapperspb.Int64Value
	if !in.IsNull() && !in.IsUnknown() {
		w = &wrapperspb.Int64Value{
			Value: in.ValueInt64(),
		}
	}

	return w
}

func ExpandAccess[V any, T accessModel[V]](ctx context.Context, cfgAccess types.Object, diags *diag.Diagnostics) T {
	var access Access
	diags.Append(cfgAccess.As(ctx, &access, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)
	if diags.HasError() {
		return nil
	}
	ac := T(new(V))
	ac.SetDataLens(access.DataLens.ValueBool())
	ac.SetDataTransfer(access.DataTransfer.ValueBool())
	ac.SetServerless(access.Serverless.ValueBool())
	ac.SetWebSql(access.WebSql.ValueBool())
	return ac
}
