package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookFunctionsRequestBuilder provides operations to manage the functions property of the microsoft.graph.workbook entity.
type ItemItemsItemWorkbookFunctionsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookFunctionsRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookFunctionsRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemItemsItemWorkbookFunctionsRequestBuilderGetQueryParameters get functions from drives
type ItemItemsItemWorkbookFunctionsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemItemsItemWorkbookFunctionsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookFunctionsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemItemsItemWorkbookFunctionsRequestBuilderGetQueryParameters
}
// ItemItemsItemWorkbookFunctionsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookFunctionsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Abs provides operations to call the abs method.
// returns a *ItemItemsItemWorkbookFunctionsAbsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Abs()(*ItemItemsItemWorkbookFunctionsAbsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAbsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AccrInt provides operations to call the accrInt method.
// returns a *ItemItemsItemWorkbookFunctionsAccrIntRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AccrInt()(*ItemItemsItemWorkbookFunctionsAccrIntRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAccrIntRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AccrIntM provides operations to call the accrIntM method.
// returns a *ItemItemsItemWorkbookFunctionsAccrIntMRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AccrIntM()(*ItemItemsItemWorkbookFunctionsAccrIntMRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAccrIntMRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Acos provides operations to call the acos method.
// returns a *ItemItemsItemWorkbookFunctionsAcosRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Acos()(*ItemItemsItemWorkbookFunctionsAcosRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAcosRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Acosh provides operations to call the acosh method.
// returns a *ItemItemsItemWorkbookFunctionsAcoshRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Acosh()(*ItemItemsItemWorkbookFunctionsAcoshRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAcoshRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Acot provides operations to call the acot method.
// returns a *ItemItemsItemWorkbookFunctionsAcotRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Acot()(*ItemItemsItemWorkbookFunctionsAcotRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAcotRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Acoth provides operations to call the acoth method.
// returns a *ItemItemsItemWorkbookFunctionsAcothRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Acoth()(*ItemItemsItemWorkbookFunctionsAcothRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAcothRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AmorDegrc provides operations to call the amorDegrc method.
// returns a *ItemItemsItemWorkbookFunctionsAmorDegrcRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AmorDegrc()(*ItemItemsItemWorkbookFunctionsAmorDegrcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAmorDegrcRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AmorLinc provides operations to call the amorLinc method.
// returns a *ItemItemsItemWorkbookFunctionsAmorLincRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AmorLinc()(*ItemItemsItemWorkbookFunctionsAmorLincRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAmorLincRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// And provides operations to call the and method.
// returns a *ItemItemsItemWorkbookFunctionsAndRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) And()(*ItemItemsItemWorkbookFunctionsAndRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAndRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Arabic provides operations to call the arabic method.
// returns a *ItemItemsItemWorkbookFunctionsArabicRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Arabic()(*ItemItemsItemWorkbookFunctionsArabicRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsArabicRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Areas provides operations to call the areas method.
// returns a *ItemItemsItemWorkbookFunctionsAreasRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Areas()(*ItemItemsItemWorkbookFunctionsAreasRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAreasRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Asc provides operations to call the asc method.
// returns a *ItemItemsItemWorkbookFunctionsAscRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Asc()(*ItemItemsItemWorkbookFunctionsAscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAscRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Asin provides operations to call the asin method.
// returns a *ItemItemsItemWorkbookFunctionsAsinRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Asin()(*ItemItemsItemWorkbookFunctionsAsinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAsinRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Asinh provides operations to call the asinh method.
// returns a *ItemItemsItemWorkbookFunctionsAsinhRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Asinh()(*ItemItemsItemWorkbookFunctionsAsinhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAsinhRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Atan provides operations to call the atan method.
// returns a *ItemItemsItemWorkbookFunctionsAtanRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Atan()(*ItemItemsItemWorkbookFunctionsAtanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAtanRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Atan2 provides operations to call the atan2 method.
// returns a *ItemItemsItemWorkbookFunctionsAtan2RequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Atan2()(*ItemItemsItemWorkbookFunctionsAtan2RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAtan2RequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Atanh provides operations to call the atanh method.
// returns a *ItemItemsItemWorkbookFunctionsAtanhRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Atanh()(*ItemItemsItemWorkbookFunctionsAtanhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAtanhRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AveDev provides operations to call the aveDev method.
// returns a *ItemItemsItemWorkbookFunctionsAveDevRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AveDev()(*ItemItemsItemWorkbookFunctionsAveDevRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAveDevRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Average provides operations to call the average method.
// returns a *ItemItemsItemWorkbookFunctionsAverageRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Average()(*ItemItemsItemWorkbookFunctionsAverageRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAverageRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AverageA provides operations to call the averageA method.
// returns a *ItemItemsItemWorkbookFunctionsAverageARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AverageA()(*ItemItemsItemWorkbookFunctionsAverageARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAverageARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AverageIf provides operations to call the averageIf method.
// returns a *ItemItemsItemWorkbookFunctionsAverageIfRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AverageIf()(*ItemItemsItemWorkbookFunctionsAverageIfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAverageIfRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AverageIfs provides operations to call the averageIfs method.
// returns a *ItemItemsItemWorkbookFunctionsAverageIfsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) AverageIfs()(*ItemItemsItemWorkbookFunctionsAverageIfsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsAverageIfsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// BahtText provides operations to call the bahtText method.
// returns a *ItemItemsItemWorkbookFunctionsBahtTextRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BahtText()(*ItemItemsItemWorkbookFunctionsBahtTextRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBahtTextRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Base provides operations to call the base method.
// returns a *ItemItemsItemWorkbookFunctionsBaseRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Base()(*ItemItemsItemWorkbookFunctionsBaseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBaseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// BesselI provides operations to call the besselI method.
// returns a *ItemItemsItemWorkbookFunctionsBesselIRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BesselI()(*ItemItemsItemWorkbookFunctionsBesselIRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBesselIRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// BesselJ provides operations to call the besselJ method.
// returns a *ItemItemsItemWorkbookFunctionsBesselJRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BesselJ()(*ItemItemsItemWorkbookFunctionsBesselJRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBesselJRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// BesselK provides operations to call the besselK method.
// returns a *ItemItemsItemWorkbookFunctionsBesselKRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BesselK()(*ItemItemsItemWorkbookFunctionsBesselKRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBesselKRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// BesselY provides operations to call the besselY method.
// returns a *ItemItemsItemWorkbookFunctionsBesselYRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) BesselY()(*ItemItemsItemWorkbookFunctionsBesselYRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBesselYRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Beta_Dist provides operations to call the beta_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsBeta_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Beta_Dist()(*ItemItemsItemWorkbookFunctionsBeta_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBeta_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Beta_Inv provides operations to call the beta_Inv method.
// returns a *ItemItemsItemWorkbookFunctionsBeta_InvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Beta_Inv()(*ItemItemsItemWorkbookFunctionsBeta_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBeta_InvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Bin2Dec provides operations to call the bin2Dec method.
// returns a *ItemItemsItemWorkbookFunctionsBin2DecRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bin2Dec()(*ItemItemsItemWorkbookFunctionsBin2DecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBin2DecRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Bin2Hex provides operations to call the bin2Hex method.
// returns a *ItemItemsItemWorkbookFunctionsBin2HexRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bin2Hex()(*ItemItemsItemWorkbookFunctionsBin2HexRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBin2HexRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Bin2Oct provides operations to call the bin2Oct method.
// returns a *ItemItemsItemWorkbookFunctionsBin2OctRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bin2Oct()(*ItemItemsItemWorkbookFunctionsBin2OctRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBin2OctRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Binom_Dist provides operations to call the binom_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsBinom_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Binom_Dist()(*ItemItemsItemWorkbookFunctionsBinom_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBinom_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Binom_Dist_Range provides operations to call the binom_Dist_Range method.
// returns a *ItemItemsItemWorkbookFunctionsBinom_Dist_RangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Binom_Dist_Range()(*ItemItemsItemWorkbookFunctionsBinom_Dist_RangeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBinom_Dist_RangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Binom_Inv provides operations to call the binom_Inv method.
// returns a *ItemItemsItemWorkbookFunctionsBinom_InvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Binom_Inv()(*ItemItemsItemWorkbookFunctionsBinom_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBinom_InvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Bitand provides operations to call the bitand method.
// returns a *ItemItemsItemWorkbookFunctionsBitandRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitand()(*ItemItemsItemWorkbookFunctionsBitandRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitandRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Bitlshift provides operations to call the bitlshift method.
// returns a *ItemItemsItemWorkbookFunctionsBitlshiftRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitlshift()(*ItemItemsItemWorkbookFunctionsBitlshiftRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitlshiftRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Bitor provides operations to call the bitor method.
// returns a *ItemItemsItemWorkbookFunctionsBitorRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitor()(*ItemItemsItemWorkbookFunctionsBitorRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitorRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Bitrshift provides operations to call the bitrshift method.
// returns a *ItemItemsItemWorkbookFunctionsBitrshiftRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitrshift()(*ItemItemsItemWorkbookFunctionsBitrshiftRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitrshiftRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Bitxor provides operations to call the bitxor method.
// returns a *ItemItemsItemWorkbookFunctionsBitxorRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Bitxor()(*ItemItemsItemWorkbookFunctionsBitxorRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsBitxorRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ceiling_Math provides operations to call the ceiling_Math method.
// returns a *ItemItemsItemWorkbookFunctionsCeiling_MathRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ceiling_Math()(*ItemItemsItemWorkbookFunctionsCeiling_MathRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCeiling_MathRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ceiling_Precise provides operations to call the ceiling_Precise method.
// returns a *ItemItemsItemWorkbookFunctionsCeiling_PreciseRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ceiling_Precise()(*ItemItemsItemWorkbookFunctionsCeiling_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCeiling_PreciseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Char provides operations to call the char method.
// returns a *ItemItemsItemWorkbookFunctionsCharRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Char()(*ItemItemsItemWorkbookFunctionsCharRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCharRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ChiSq_Dist provides operations to call the chiSq_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsChiSq_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ChiSq_Dist()(*ItemItemsItemWorkbookFunctionsChiSq_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChiSq_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ChiSq_Dist_RT provides operations to call the chiSq_Dist_RT method.
// returns a *ItemItemsItemWorkbookFunctionsChiSq_Dist_RTRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ChiSq_Dist_RT()(*ItemItemsItemWorkbookFunctionsChiSq_Dist_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChiSq_Dist_RTRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ChiSq_Inv provides operations to call the chiSq_Inv method.
// returns a *ItemItemsItemWorkbookFunctionsChiSq_InvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ChiSq_Inv()(*ItemItemsItemWorkbookFunctionsChiSq_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChiSq_InvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ChiSq_Inv_RT provides operations to call the chiSq_Inv_RT method.
// returns a *ItemItemsItemWorkbookFunctionsChiSq_Inv_RTRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ChiSq_Inv_RT()(*ItemItemsItemWorkbookFunctionsChiSq_Inv_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChiSq_Inv_RTRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Choose provides operations to call the choose method.
// returns a *ItemItemsItemWorkbookFunctionsChooseRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Choose()(*ItemItemsItemWorkbookFunctionsChooseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsChooseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Clean provides operations to call the clean method.
// returns a *ItemItemsItemWorkbookFunctionsCleanRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Clean()(*ItemItemsItemWorkbookFunctionsCleanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCleanRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Code provides operations to call the code method.
// returns a *ItemItemsItemWorkbookFunctionsCodeRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Code()(*ItemItemsItemWorkbookFunctionsCodeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCodeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Columns provides operations to call the columns method.
// returns a *ItemItemsItemWorkbookFunctionsColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Columns()(*ItemItemsItemWorkbookFunctionsColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Combin provides operations to call the combin method.
// returns a *ItemItemsItemWorkbookFunctionsCombinRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Combin()(*ItemItemsItemWorkbookFunctionsCombinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCombinRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Combina provides operations to call the combina method.
// returns a *ItemItemsItemWorkbookFunctionsCombinaRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Combina()(*ItemItemsItemWorkbookFunctionsCombinaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCombinaRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Complex provides operations to call the complex method.
// returns a *ItemItemsItemWorkbookFunctionsComplexRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Complex()(*ItemItemsItemWorkbookFunctionsComplexRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsComplexRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Concatenate provides operations to call the concatenate method.
// returns a *ItemItemsItemWorkbookFunctionsConcatenateRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Concatenate()(*ItemItemsItemWorkbookFunctionsConcatenateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsConcatenateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Confidence_Norm provides operations to call the confidence_Norm method.
// returns a *ItemItemsItemWorkbookFunctionsConfidence_NormRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Confidence_Norm()(*ItemItemsItemWorkbookFunctionsConfidence_NormRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsConfidence_NormRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Confidence_T provides operations to call the confidence_T method.
// returns a *ItemItemsItemWorkbookFunctionsConfidence_TRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Confidence_T()(*ItemItemsItemWorkbookFunctionsConfidence_TRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsConfidence_TRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemItemsItemWorkbookFunctionsRequestBuilderInternal instantiates a new ItemItemsItemWorkbookFunctionsRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookFunctionsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookFunctionsRequestBuilder) {
    m := &ItemItemsItemWorkbookFunctionsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/functions{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookFunctionsRequestBuilder instantiates a new ItemItemsItemWorkbookFunctionsRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookFunctionsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookFunctionsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookFunctionsRequestBuilderInternal(urlParams, requestAdapter)
}
// Convert provides operations to call the convert method.
// returns a *ItemItemsItemWorkbookFunctionsConvertRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Convert()(*ItemItemsItemWorkbookFunctionsConvertRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsConvertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Cos provides operations to call the cos method.
// returns a *ItemItemsItemWorkbookFunctionsCosRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Cos()(*ItemItemsItemWorkbookFunctionsCosRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCosRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Cosh provides operations to call the cosh method.
// returns a *ItemItemsItemWorkbookFunctionsCoshRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Cosh()(*ItemItemsItemWorkbookFunctionsCoshRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoshRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Cot provides operations to call the cot method.
// returns a *ItemItemsItemWorkbookFunctionsCotRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Cot()(*ItemItemsItemWorkbookFunctionsCotRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCotRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Coth provides operations to call the coth method.
// returns a *ItemItemsItemWorkbookFunctionsCothRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Coth()(*ItemItemsItemWorkbookFunctionsCothRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCothRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Count provides operations to call the count method.
// returns a *ItemItemsItemWorkbookFunctionsCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Count()(*ItemItemsItemWorkbookFunctionsCountRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CountA provides operations to call the countA method.
// returns a *ItemItemsItemWorkbookFunctionsCountARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CountA()(*ItemItemsItemWorkbookFunctionsCountARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CountBlank provides operations to call the countBlank method.
// returns a *ItemItemsItemWorkbookFunctionsCountBlankRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CountBlank()(*ItemItemsItemWorkbookFunctionsCountBlankRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountBlankRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CountIf provides operations to call the countIf method.
// returns a *ItemItemsItemWorkbookFunctionsCountIfRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CountIf()(*ItemItemsItemWorkbookFunctionsCountIfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountIfRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CountIfs provides operations to call the countIfs method.
// returns a *ItemItemsItemWorkbookFunctionsCountIfsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CountIfs()(*ItemItemsItemWorkbookFunctionsCountIfsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCountIfsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CoupDayBs provides operations to call the coupDayBs method.
// returns a *ItemItemsItemWorkbookFunctionsCoupDayBsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupDayBs()(*ItemItemsItemWorkbookFunctionsCoupDayBsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupDayBsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CoupDays provides operations to call the coupDays method.
// returns a *ItemItemsItemWorkbookFunctionsCoupDaysRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupDays()(*ItemItemsItemWorkbookFunctionsCoupDaysRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupDaysRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CoupDaysNc provides operations to call the coupDaysNc method.
// returns a *ItemItemsItemWorkbookFunctionsCoupDaysNcRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupDaysNc()(*ItemItemsItemWorkbookFunctionsCoupDaysNcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupDaysNcRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CoupNcd provides operations to call the coupNcd method.
// returns a *ItemItemsItemWorkbookFunctionsCoupNcdRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupNcd()(*ItemItemsItemWorkbookFunctionsCoupNcdRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupNcdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CoupNum provides operations to call the coupNum method.
// returns a *ItemItemsItemWorkbookFunctionsCoupNumRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupNum()(*ItemItemsItemWorkbookFunctionsCoupNumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupNumRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CoupPcd provides operations to call the coupPcd method.
// returns a *ItemItemsItemWorkbookFunctionsCoupPcdRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CoupPcd()(*ItemItemsItemWorkbookFunctionsCoupPcdRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCoupPcdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Csc provides operations to call the csc method.
// returns a *ItemItemsItemWorkbookFunctionsCscRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Csc()(*ItemItemsItemWorkbookFunctionsCscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCscRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Csch provides operations to call the csch method.
// returns a *ItemItemsItemWorkbookFunctionsCschRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Csch()(*ItemItemsItemWorkbookFunctionsCschRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCschRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CumIPmt provides operations to call the cumIPmt method.
// returns a *ItemItemsItemWorkbookFunctionsCumIPmtRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CumIPmt()(*ItemItemsItemWorkbookFunctionsCumIPmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCumIPmtRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CumPrinc provides operations to call the cumPrinc method.
// returns a *ItemItemsItemWorkbookFunctionsCumPrincRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) CumPrinc()(*ItemItemsItemWorkbookFunctionsCumPrincRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsCumPrincRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Date provides operations to call the date method.
// returns a *ItemItemsItemWorkbookFunctionsDateRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Date()(*ItemItemsItemWorkbookFunctionsDateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Datevalue provides operations to call the datevalue method.
// returns a *ItemItemsItemWorkbookFunctionsDatevalueRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Datevalue()(*ItemItemsItemWorkbookFunctionsDatevalueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDatevalueRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Daverage provides operations to call the daverage method.
// returns a *ItemItemsItemWorkbookFunctionsDaverageRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Daverage()(*ItemItemsItemWorkbookFunctionsDaverageRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDaverageRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Day provides operations to call the day method.
// returns a *ItemItemsItemWorkbookFunctionsDayRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Day()(*ItemItemsItemWorkbookFunctionsDayRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDayRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Days provides operations to call the days method.
// returns a *ItemItemsItemWorkbookFunctionsDaysRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Days()(*ItemItemsItemWorkbookFunctionsDaysRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDaysRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Days360 provides operations to call the days360 method.
// returns a *ItemItemsItemWorkbookFunctionsDays360RequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Days360()(*ItemItemsItemWorkbookFunctionsDays360RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDays360RequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Db provides operations to call the db method.
// returns a *ItemItemsItemWorkbookFunctionsDbRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Db()(*ItemItemsItemWorkbookFunctionsDbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDbRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dbcs provides operations to call the dbcs method.
// returns a *ItemItemsItemWorkbookFunctionsDbcsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dbcs()(*ItemItemsItemWorkbookFunctionsDbcsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDbcsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dcount provides operations to call the dcount method.
// returns a *ItemItemsItemWorkbookFunctionsDcountRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dcount()(*ItemItemsItemWorkbookFunctionsDcountRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDcountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DcountA provides operations to call the dcountA method.
// returns a *ItemItemsItemWorkbookFunctionsDcountARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DcountA()(*ItemItemsItemWorkbookFunctionsDcountARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDcountARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ddb provides operations to call the ddb method.
// returns a *ItemItemsItemWorkbookFunctionsDdbRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ddb()(*ItemItemsItemWorkbookFunctionsDdbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDdbRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dec2Bin provides operations to call the dec2Bin method.
// returns a *ItemItemsItemWorkbookFunctionsDec2BinRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dec2Bin()(*ItemItemsItemWorkbookFunctionsDec2BinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDec2BinRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dec2Hex provides operations to call the dec2Hex method.
// returns a *ItemItemsItemWorkbookFunctionsDec2HexRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dec2Hex()(*ItemItemsItemWorkbookFunctionsDec2HexRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDec2HexRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dec2Oct provides operations to call the dec2Oct method.
// returns a *ItemItemsItemWorkbookFunctionsDec2OctRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dec2Oct()(*ItemItemsItemWorkbookFunctionsDec2OctRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDec2OctRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Decimal provides operations to call the decimal method.
// returns a *ItemItemsItemWorkbookFunctionsDecimalRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Decimal()(*ItemItemsItemWorkbookFunctionsDecimalRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDecimalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Degrees provides operations to call the degrees method.
// returns a *ItemItemsItemWorkbookFunctionsDegreesRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Degrees()(*ItemItemsItemWorkbookFunctionsDegreesRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDegreesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Delete delete navigation property functions for drives
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Delta provides operations to call the delta method.
// returns a *ItemItemsItemWorkbookFunctionsDeltaRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Delta()(*ItemItemsItemWorkbookFunctionsDeltaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDeltaRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DevSq provides operations to call the devSq method.
// returns a *ItemItemsItemWorkbookFunctionsDevSqRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DevSq()(*ItemItemsItemWorkbookFunctionsDevSqRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDevSqRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dget provides operations to call the dget method.
// returns a *ItemItemsItemWorkbookFunctionsDgetRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dget()(*ItemItemsItemWorkbookFunctionsDgetRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDgetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Disc provides operations to call the disc method.
// returns a *ItemItemsItemWorkbookFunctionsDiscRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Disc()(*ItemItemsItemWorkbookFunctionsDiscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDiscRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dmax provides operations to call the dmax method.
// returns a *ItemItemsItemWorkbookFunctionsDmaxRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dmax()(*ItemItemsItemWorkbookFunctionsDmaxRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDmaxRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dmin provides operations to call the dmin method.
// returns a *ItemItemsItemWorkbookFunctionsDminRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dmin()(*ItemItemsItemWorkbookFunctionsDminRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDminRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dollar provides operations to call the dollar method.
// returns a *ItemItemsItemWorkbookFunctionsDollarRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dollar()(*ItemItemsItemWorkbookFunctionsDollarRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDollarRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DollarDe provides operations to call the dollarDe method.
// returns a *ItemItemsItemWorkbookFunctionsDollarDeRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DollarDe()(*ItemItemsItemWorkbookFunctionsDollarDeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDollarDeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DollarFr provides operations to call the dollarFr method.
// returns a *ItemItemsItemWorkbookFunctionsDollarFrRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DollarFr()(*ItemItemsItemWorkbookFunctionsDollarFrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDollarFrRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dproduct provides operations to call the dproduct method.
// returns a *ItemItemsItemWorkbookFunctionsDproductRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dproduct()(*ItemItemsItemWorkbookFunctionsDproductRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDproductRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DstDev provides operations to call the dstDev method.
// returns a *ItemItemsItemWorkbookFunctionsDstDevRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DstDev()(*ItemItemsItemWorkbookFunctionsDstDevRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDstDevRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DstDevP provides operations to call the dstDevP method.
// returns a *ItemItemsItemWorkbookFunctionsDstDevPRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DstDevP()(*ItemItemsItemWorkbookFunctionsDstDevPRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDstDevPRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dsum provides operations to call the dsum method.
// returns a *ItemItemsItemWorkbookFunctionsDsumRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dsum()(*ItemItemsItemWorkbookFunctionsDsumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDsumRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Duration provides operations to call the duration method.
// returns a *ItemItemsItemWorkbookFunctionsDurationRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Duration()(*ItemItemsItemWorkbookFunctionsDurationRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDurationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Dvar provides operations to call the dvar method.
// returns a *ItemItemsItemWorkbookFunctionsDvarRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Dvar()(*ItemItemsItemWorkbookFunctionsDvarRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDvarRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DvarP provides operations to call the dvarP method.
// returns a *ItemItemsItemWorkbookFunctionsDvarPRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) DvarP()(*ItemItemsItemWorkbookFunctionsDvarPRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsDvarPRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ecma_Ceiling provides operations to call the ecma_Ceiling method.
// returns a *ItemItemsItemWorkbookFunctionsEcma_CeilingRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ecma_Ceiling()(*ItemItemsItemWorkbookFunctionsEcma_CeilingRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEcma_CeilingRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Edate provides operations to call the edate method.
// returns a *ItemItemsItemWorkbookFunctionsEdateRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Edate()(*ItemItemsItemWorkbookFunctionsEdateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEdateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Effect provides operations to call the effect method.
// returns a *ItemItemsItemWorkbookFunctionsEffectRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Effect()(*ItemItemsItemWorkbookFunctionsEffectRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEffectRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EoMonth provides operations to call the eoMonth method.
// returns a *ItemItemsItemWorkbookFunctionsEoMonthRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) EoMonth()(*ItemItemsItemWorkbookFunctionsEoMonthRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEoMonthRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Erf provides operations to call the erf method.
// returns a *ItemItemsItemWorkbookFunctionsErfRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Erf()(*ItemItemsItemWorkbookFunctionsErfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsErfRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Erf_Precise provides operations to call the erf_Precise method.
// returns a *ItemItemsItemWorkbookFunctionsErf_PreciseRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Erf_Precise()(*ItemItemsItemWorkbookFunctionsErf_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsErf_PreciseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ErfC provides operations to call the erfC method.
// returns a *ItemItemsItemWorkbookFunctionsErfCRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ErfC()(*ItemItemsItemWorkbookFunctionsErfCRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsErfCRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ErfC_Precise provides operations to call the erfC_Precise method.
// returns a *ItemItemsItemWorkbookFunctionsErfC_PreciseRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ErfC_Precise()(*ItemItemsItemWorkbookFunctionsErfC_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsErfC_PreciseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Error_Type provides operations to call the error_Type method.
// returns a *ItemItemsItemWorkbookFunctionsError_TypeRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Error_Type()(*ItemItemsItemWorkbookFunctionsError_TypeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsError_TypeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Even provides operations to call the even method.
// returns a *ItemItemsItemWorkbookFunctionsEvenRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Even()(*ItemItemsItemWorkbookFunctionsEvenRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsEvenRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Exact provides operations to call the exact method.
// returns a *ItemItemsItemWorkbookFunctionsExactRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Exact()(*ItemItemsItemWorkbookFunctionsExactRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsExactRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Exp provides operations to call the exp method.
// returns a *ItemItemsItemWorkbookFunctionsExpRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Exp()(*ItemItemsItemWorkbookFunctionsExpRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsExpRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Expon_Dist provides operations to call the expon_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsExpon_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Expon_Dist()(*ItemItemsItemWorkbookFunctionsExpon_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsExpon_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// F_Dist provides operations to call the f_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsF_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) F_Dist()(*ItemItemsItemWorkbookFunctionsF_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsF_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// F_Dist_RT provides operations to call the f_Dist_RT method.
// returns a *ItemItemsItemWorkbookFunctionsF_Dist_RTRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) F_Dist_RT()(*ItemItemsItemWorkbookFunctionsF_Dist_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsF_Dist_RTRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// F_Inv provides operations to call the f_Inv method.
// returns a *ItemItemsItemWorkbookFunctionsF_InvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) F_Inv()(*ItemItemsItemWorkbookFunctionsF_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsF_InvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// F_Inv_RT provides operations to call the f_Inv_RT method.
// returns a *ItemItemsItemWorkbookFunctionsF_Inv_RTRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) F_Inv_RT()(*ItemItemsItemWorkbookFunctionsF_Inv_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsF_Inv_RTRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Fact provides operations to call the fact method.
// returns a *ItemItemsItemWorkbookFunctionsFactRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fact()(*ItemItemsItemWorkbookFunctionsFactRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFactRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// FactDouble provides operations to call the factDouble method.
// returns a *ItemItemsItemWorkbookFunctionsFactDoubleRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) FactDouble()(*ItemItemsItemWorkbookFunctionsFactDoubleRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFactDoubleRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// False provides operations to call the false method.
// returns a *ItemItemsItemWorkbookFunctionsFalseRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) False()(*ItemItemsItemWorkbookFunctionsFalseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFalseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Find provides operations to call the find method.
// returns a *ItemItemsItemWorkbookFunctionsFindRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Find()(*ItemItemsItemWorkbookFunctionsFindRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFindRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// FindB provides operations to call the findB method.
// returns a *ItemItemsItemWorkbookFunctionsFindBRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) FindB()(*ItemItemsItemWorkbookFunctionsFindBRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFindBRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Fisher provides operations to call the fisher method.
// returns a *ItemItemsItemWorkbookFunctionsFisherRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fisher()(*ItemItemsItemWorkbookFunctionsFisherRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFisherRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// FisherInv provides operations to call the fisherInv method.
// returns a *ItemItemsItemWorkbookFunctionsFisherInvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) FisherInv()(*ItemItemsItemWorkbookFunctionsFisherInvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFisherInvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Fixed provides operations to call the fixed method.
// returns a *ItemItemsItemWorkbookFunctionsFixedRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fixed()(*ItemItemsItemWorkbookFunctionsFixedRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFixedRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Floor_Math provides operations to call the floor_Math method.
// returns a *ItemItemsItemWorkbookFunctionsFloor_MathRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Floor_Math()(*ItemItemsItemWorkbookFunctionsFloor_MathRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFloor_MathRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Floor_Precise provides operations to call the floor_Precise method.
// returns a *ItemItemsItemWorkbookFunctionsFloor_PreciseRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Floor_Precise()(*ItemItemsItemWorkbookFunctionsFloor_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFloor_PreciseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Fv provides operations to call the fv method.
// returns a *ItemItemsItemWorkbookFunctionsFvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fv()(*ItemItemsItemWorkbookFunctionsFvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Fvschedule provides operations to call the fvschedule method.
// returns a *ItemItemsItemWorkbookFunctionsFvscheduleRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Fvschedule()(*ItemItemsItemWorkbookFunctionsFvscheduleRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsFvscheduleRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Gamma provides operations to call the gamma method.
// returns a *ItemItemsItemWorkbookFunctionsGammaRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gamma()(*ItemItemsItemWorkbookFunctionsGammaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGammaRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Gamma_Dist provides operations to call the gamma_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsGamma_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gamma_Dist()(*ItemItemsItemWorkbookFunctionsGamma_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGamma_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Gamma_Inv provides operations to call the gamma_Inv method.
// returns a *ItemItemsItemWorkbookFunctionsGamma_InvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gamma_Inv()(*ItemItemsItemWorkbookFunctionsGamma_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGamma_InvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GammaLn provides operations to call the gammaLn method.
// returns a *ItemItemsItemWorkbookFunctionsGammaLnRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) GammaLn()(*ItemItemsItemWorkbookFunctionsGammaLnRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGammaLnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GammaLn_Precise provides operations to call the gammaLn_Precise method.
// returns a *ItemItemsItemWorkbookFunctionsGammaLn_PreciseRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) GammaLn_Precise()(*ItemItemsItemWorkbookFunctionsGammaLn_PreciseRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGammaLn_PreciseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Gauss provides operations to call the gauss method.
// returns a *ItemItemsItemWorkbookFunctionsGaussRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gauss()(*ItemItemsItemWorkbookFunctionsGaussRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGaussRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Gcd provides operations to call the gcd method.
// returns a *ItemItemsItemWorkbookFunctionsGcdRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Gcd()(*ItemItemsItemWorkbookFunctionsGcdRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGcdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GeoMean provides operations to call the geoMean method.
// returns a *ItemItemsItemWorkbookFunctionsGeoMeanRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) GeoMean()(*ItemItemsItemWorkbookFunctionsGeoMeanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGeoMeanRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GeStep provides operations to call the geStep method.
// returns a *ItemItemsItemWorkbookFunctionsGeStepRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) GeStep()(*ItemItemsItemWorkbookFunctionsGeStepRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsGeStepRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get functions from drives
// returns a WorkbookFunctionsable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWorkbookFunctionsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable), nil
}
// HarMean provides operations to call the harMean method.
// returns a *ItemItemsItemWorkbookFunctionsHarMeanRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) HarMean()(*ItemItemsItemWorkbookFunctionsHarMeanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHarMeanRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Hex2Bin provides operations to call the hex2Bin method.
// returns a *ItemItemsItemWorkbookFunctionsHex2BinRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hex2Bin()(*ItemItemsItemWorkbookFunctionsHex2BinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHex2BinRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Hex2Dec provides operations to call the hex2Dec method.
// returns a *ItemItemsItemWorkbookFunctionsHex2DecRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hex2Dec()(*ItemItemsItemWorkbookFunctionsHex2DecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHex2DecRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Hex2Oct provides operations to call the hex2Oct method.
// returns a *ItemItemsItemWorkbookFunctionsHex2OctRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hex2Oct()(*ItemItemsItemWorkbookFunctionsHex2OctRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHex2OctRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Hlookup provides operations to call the hlookup method.
// returns a *ItemItemsItemWorkbookFunctionsHlookupRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hlookup()(*ItemItemsItemWorkbookFunctionsHlookupRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHlookupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Hour provides operations to call the hour method.
// returns a *ItemItemsItemWorkbookFunctionsHourRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hour()(*ItemItemsItemWorkbookFunctionsHourRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHourRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Hyperlink provides operations to call the hyperlink method.
// returns a *ItemItemsItemWorkbookFunctionsHyperlinkRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Hyperlink()(*ItemItemsItemWorkbookFunctionsHyperlinkRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHyperlinkRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// HypGeom_Dist provides operations to call the hypGeom_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsHypGeom_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) HypGeom_Dist()(*ItemItemsItemWorkbookFunctionsHypGeom_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsHypGeom_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IfEscaped provides operations to call the if method.
// returns a *ItemItemsItemWorkbookFunctionsIfRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IfEscaped()(*ItemItemsItemWorkbookFunctionsIfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIfRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImAbs provides operations to call the imAbs method.
// returns a *ItemItemsItemWorkbookFunctionsImAbsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImAbs()(*ItemItemsItemWorkbookFunctionsImAbsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImAbsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Imaginary provides operations to call the imaginary method.
// returns a *ItemItemsItemWorkbookFunctionsImaginaryRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Imaginary()(*ItemItemsItemWorkbookFunctionsImaginaryRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImaginaryRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImArgument provides operations to call the imArgument method.
// returns a *ItemItemsItemWorkbookFunctionsImArgumentRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImArgument()(*ItemItemsItemWorkbookFunctionsImArgumentRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImArgumentRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImConjugate provides operations to call the imConjugate method.
// returns a *ItemItemsItemWorkbookFunctionsImConjugateRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImConjugate()(*ItemItemsItemWorkbookFunctionsImConjugateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImConjugateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImCos provides operations to call the imCos method.
// returns a *ItemItemsItemWorkbookFunctionsImCosRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCos()(*ItemItemsItemWorkbookFunctionsImCosRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCosRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImCosh provides operations to call the imCosh method.
// returns a *ItemItemsItemWorkbookFunctionsImCoshRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCosh()(*ItemItemsItemWorkbookFunctionsImCoshRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCoshRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImCot provides operations to call the imCot method.
// returns a *ItemItemsItemWorkbookFunctionsImCotRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCot()(*ItemItemsItemWorkbookFunctionsImCotRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCotRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImCsc provides operations to call the imCsc method.
// returns a *ItemItemsItemWorkbookFunctionsImCscRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCsc()(*ItemItemsItemWorkbookFunctionsImCscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCscRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImCsch provides operations to call the imCsch method.
// returns a *ItemItemsItemWorkbookFunctionsImCschRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImCsch()(*ItemItemsItemWorkbookFunctionsImCschRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImCschRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImDiv provides operations to call the imDiv method.
// returns a *ItemItemsItemWorkbookFunctionsImDivRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImDiv()(*ItemItemsItemWorkbookFunctionsImDivRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImDivRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImExp provides operations to call the imExp method.
// returns a *ItemItemsItemWorkbookFunctionsImExpRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImExp()(*ItemItemsItemWorkbookFunctionsImExpRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImExpRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImLn provides operations to call the imLn method.
// returns a *ItemItemsItemWorkbookFunctionsImLnRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImLn()(*ItemItemsItemWorkbookFunctionsImLnRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImLnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImLog10 provides operations to call the imLog10 method.
// returns a *ItemItemsItemWorkbookFunctionsImLog10RequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImLog10()(*ItemItemsItemWorkbookFunctionsImLog10RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImLog10RequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImLog2 provides operations to call the imLog2 method.
// returns a *ItemItemsItemWorkbookFunctionsImLog2RequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImLog2()(*ItemItemsItemWorkbookFunctionsImLog2RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImLog2RequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImPower provides operations to call the imPower method.
// returns a *ItemItemsItemWorkbookFunctionsImPowerRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImPower()(*ItemItemsItemWorkbookFunctionsImPowerRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImPowerRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImProduct provides operations to call the imProduct method.
// returns a *ItemItemsItemWorkbookFunctionsImProductRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImProduct()(*ItemItemsItemWorkbookFunctionsImProductRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImProductRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImReal provides operations to call the imReal method.
// returns a *ItemItemsItemWorkbookFunctionsImRealRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImReal()(*ItemItemsItemWorkbookFunctionsImRealRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImRealRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImSec provides operations to call the imSec method.
// returns a *ItemItemsItemWorkbookFunctionsImSecRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSec()(*ItemItemsItemWorkbookFunctionsImSecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSecRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImSech provides operations to call the imSech method.
// returns a *ItemItemsItemWorkbookFunctionsImSechRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSech()(*ItemItemsItemWorkbookFunctionsImSechRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSechRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImSin provides operations to call the imSin method.
// returns a *ItemItemsItemWorkbookFunctionsImSinRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSin()(*ItemItemsItemWorkbookFunctionsImSinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSinRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImSinh provides operations to call the imSinh method.
// returns a *ItemItemsItemWorkbookFunctionsImSinhRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSinh()(*ItemItemsItemWorkbookFunctionsImSinhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSinhRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImSqrt provides operations to call the imSqrt method.
// returns a *ItemItemsItemWorkbookFunctionsImSqrtRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSqrt()(*ItemItemsItemWorkbookFunctionsImSqrtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSqrtRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImSub provides operations to call the imSub method.
// returns a *ItemItemsItemWorkbookFunctionsImSubRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSub()(*ItemItemsItemWorkbookFunctionsImSubRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSubRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImSum provides operations to call the imSum method.
// returns a *ItemItemsItemWorkbookFunctionsImSumRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImSum()(*ItemItemsItemWorkbookFunctionsImSumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImSumRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ImTan provides operations to call the imTan method.
// returns a *ItemItemsItemWorkbookFunctionsImTanRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ImTan()(*ItemItemsItemWorkbookFunctionsImTanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsImTanRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Int provides operations to call the int method.
// returns a *ItemItemsItemWorkbookFunctionsIntRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Int()(*ItemItemsItemWorkbookFunctionsIntRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIntRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntRate provides operations to call the intRate method.
// returns a *ItemItemsItemWorkbookFunctionsIntRateRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IntRate()(*ItemItemsItemWorkbookFunctionsIntRateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIntRateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ipmt provides operations to call the ipmt method.
// returns a *ItemItemsItemWorkbookFunctionsIpmtRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ipmt()(*ItemItemsItemWorkbookFunctionsIpmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIpmtRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Irr provides operations to call the irr method.
// returns a *ItemItemsItemWorkbookFunctionsIrrRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Irr()(*ItemItemsItemWorkbookFunctionsIrrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIrrRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsErr provides operations to call the isErr method.
// returns a *ItemItemsItemWorkbookFunctionsIsErrRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsErr()(*ItemItemsItemWorkbookFunctionsIsErrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsErrRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsError provides operations to call the isError method.
// returns a *ItemItemsItemWorkbookFunctionsIsErrorRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsError()(*ItemItemsItemWorkbookFunctionsIsErrorRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsErrorRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsEven provides operations to call the isEven method.
// returns a *ItemItemsItemWorkbookFunctionsIsEvenRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsEven()(*ItemItemsItemWorkbookFunctionsIsEvenRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsEvenRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsFormula provides operations to call the isFormula method.
// returns a *ItemItemsItemWorkbookFunctionsIsFormulaRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsFormula()(*ItemItemsItemWorkbookFunctionsIsFormulaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsFormulaRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsLogical provides operations to call the isLogical method.
// returns a *ItemItemsItemWorkbookFunctionsIsLogicalRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsLogical()(*ItemItemsItemWorkbookFunctionsIsLogicalRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsLogicalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsNA provides operations to call the isNA method.
// returns a *ItemItemsItemWorkbookFunctionsIsNARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsNA()(*ItemItemsItemWorkbookFunctionsIsNARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsNARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsNonText provides operations to call the isNonText method.
// returns a *ItemItemsItemWorkbookFunctionsIsNonTextRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsNonText()(*ItemItemsItemWorkbookFunctionsIsNonTextRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsNonTextRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsNumber provides operations to call the isNumber method.
// returns a *ItemItemsItemWorkbookFunctionsIsNumberRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsNumber()(*ItemItemsItemWorkbookFunctionsIsNumberRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsNumberRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Iso_Ceiling provides operations to call the iso_Ceiling method.
// returns a *ItemItemsItemWorkbookFunctionsIso_CeilingRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Iso_Ceiling()(*ItemItemsItemWorkbookFunctionsIso_CeilingRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIso_CeilingRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsOdd provides operations to call the isOdd method.
// returns a *ItemItemsItemWorkbookFunctionsIsOddRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsOdd()(*ItemItemsItemWorkbookFunctionsIsOddRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsOddRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsoWeekNum provides operations to call the isoWeekNum method.
// returns a *ItemItemsItemWorkbookFunctionsIsoWeekNumRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsoWeekNum()(*ItemItemsItemWorkbookFunctionsIsoWeekNumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsoWeekNumRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ispmt provides operations to call the ispmt method.
// returns a *ItemItemsItemWorkbookFunctionsIspmtRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ispmt()(*ItemItemsItemWorkbookFunctionsIspmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIspmtRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Isref provides operations to call the isref method.
// returns a *ItemItemsItemWorkbookFunctionsIsrefRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Isref()(*ItemItemsItemWorkbookFunctionsIsrefRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsrefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IsText provides operations to call the isText method.
// returns a *ItemItemsItemWorkbookFunctionsIsTextRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) IsText()(*ItemItemsItemWorkbookFunctionsIsTextRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsIsTextRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Kurt provides operations to call the kurt method.
// returns a *ItemItemsItemWorkbookFunctionsKurtRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Kurt()(*ItemItemsItemWorkbookFunctionsKurtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsKurtRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Large provides operations to call the large method.
// returns a *ItemItemsItemWorkbookFunctionsLargeRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Large()(*ItemItemsItemWorkbookFunctionsLargeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLargeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Lcm provides operations to call the lcm method.
// returns a *ItemItemsItemWorkbookFunctionsLcmRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Lcm()(*ItemItemsItemWorkbookFunctionsLcmRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLcmRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Left provides operations to call the left method.
// returns a *ItemItemsItemWorkbookFunctionsLeftRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Left()(*ItemItemsItemWorkbookFunctionsLeftRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLeftRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Leftb provides operations to call the leftb method.
// returns a *ItemItemsItemWorkbookFunctionsLeftbRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Leftb()(*ItemItemsItemWorkbookFunctionsLeftbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLeftbRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Len provides operations to call the len method.
// returns a *ItemItemsItemWorkbookFunctionsLenRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Len()(*ItemItemsItemWorkbookFunctionsLenRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLenRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Lenb provides operations to call the lenb method.
// returns a *ItemItemsItemWorkbookFunctionsLenbRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Lenb()(*ItemItemsItemWorkbookFunctionsLenbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLenbRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ln provides operations to call the ln method.
// returns a *ItemItemsItemWorkbookFunctionsLnRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ln()(*ItemItemsItemWorkbookFunctionsLnRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Log provides operations to call the log method.
// returns a *ItemItemsItemWorkbookFunctionsLogRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Log()(*ItemItemsItemWorkbookFunctionsLogRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLogRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Log10 provides operations to call the log10 method.
// returns a *ItemItemsItemWorkbookFunctionsLog10RequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Log10()(*ItemItemsItemWorkbookFunctionsLog10RequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLog10RequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LogNorm_Dist provides operations to call the logNorm_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsLogNorm_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) LogNorm_Dist()(*ItemItemsItemWorkbookFunctionsLogNorm_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLogNorm_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LogNorm_Inv provides operations to call the logNorm_Inv method.
// returns a *ItemItemsItemWorkbookFunctionsLogNorm_InvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) LogNorm_Inv()(*ItemItemsItemWorkbookFunctionsLogNorm_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLogNorm_InvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Lookup provides operations to call the lookup method.
// returns a *ItemItemsItemWorkbookFunctionsLookupRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Lookup()(*ItemItemsItemWorkbookFunctionsLookupRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLookupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Lower provides operations to call the lower method.
// returns a *ItemItemsItemWorkbookFunctionsLowerRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Lower()(*ItemItemsItemWorkbookFunctionsLowerRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsLowerRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Match provides operations to call the match method.
// returns a *ItemItemsItemWorkbookFunctionsMatchRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Match()(*ItemItemsItemWorkbookFunctionsMatchRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMatchRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Max provides operations to call the max method.
// returns a *ItemItemsItemWorkbookFunctionsMaxRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Max()(*ItemItemsItemWorkbookFunctionsMaxRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMaxRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MaxA provides operations to call the maxA method.
// returns a *ItemItemsItemWorkbookFunctionsMaxARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) MaxA()(*ItemItemsItemWorkbookFunctionsMaxARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMaxARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Mduration provides operations to call the mduration method.
// returns a *ItemItemsItemWorkbookFunctionsMdurationRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mduration()(*ItemItemsItemWorkbookFunctionsMdurationRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMdurationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Median provides operations to call the median method.
// returns a *ItemItemsItemWorkbookFunctionsMedianRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Median()(*ItemItemsItemWorkbookFunctionsMedianRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMedianRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Mid provides operations to call the mid method.
// returns a *ItemItemsItemWorkbookFunctionsMidRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mid()(*ItemItemsItemWorkbookFunctionsMidRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMidRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Midb provides operations to call the midb method.
// returns a *ItemItemsItemWorkbookFunctionsMidbRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Midb()(*ItemItemsItemWorkbookFunctionsMidbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMidbRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Min provides operations to call the min method.
// returns a *ItemItemsItemWorkbookFunctionsMinRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Min()(*ItemItemsItemWorkbookFunctionsMinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMinRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MinA provides operations to call the minA method.
// returns a *ItemItemsItemWorkbookFunctionsMinARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) MinA()(*ItemItemsItemWorkbookFunctionsMinARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMinARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Minute provides operations to call the minute method.
// returns a *ItemItemsItemWorkbookFunctionsMinuteRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Minute()(*ItemItemsItemWorkbookFunctionsMinuteRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMinuteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Mirr provides operations to call the mirr method.
// returns a *ItemItemsItemWorkbookFunctionsMirrRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mirr()(*ItemItemsItemWorkbookFunctionsMirrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMirrRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Mod provides operations to call the mod method.
// returns a *ItemItemsItemWorkbookFunctionsModRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mod()(*ItemItemsItemWorkbookFunctionsModRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsModRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Month provides operations to call the month method.
// returns a *ItemItemsItemWorkbookFunctionsMonthRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Month()(*ItemItemsItemWorkbookFunctionsMonthRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMonthRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Mround provides operations to call the mround method.
// returns a *ItemItemsItemWorkbookFunctionsMroundRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Mround()(*ItemItemsItemWorkbookFunctionsMroundRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMroundRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MultiNomial provides operations to call the multiNomial method.
// returns a *ItemItemsItemWorkbookFunctionsMultiNomialRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) MultiNomial()(*ItemItemsItemWorkbookFunctionsMultiNomialRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsMultiNomialRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// N provides operations to call the n method.
// returns a *ItemItemsItemWorkbookFunctionsNRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) N()(*ItemItemsItemWorkbookFunctionsNRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Na provides operations to call the na method.
// returns a *ItemItemsItemWorkbookFunctionsNaRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Na()(*ItemItemsItemWorkbookFunctionsNaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNaRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NegBinom_Dist provides operations to call the negBinom_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsNegBinom_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) NegBinom_Dist()(*ItemItemsItemWorkbookFunctionsNegBinom_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNegBinom_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NetworkDays provides operations to call the networkDays method.
// returns a *ItemItemsItemWorkbookFunctionsNetworkDaysRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) NetworkDays()(*ItemItemsItemWorkbookFunctionsNetworkDaysRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNetworkDaysRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NetworkDays_Intl provides operations to call the networkDays_Intl method.
// returns a *ItemItemsItemWorkbookFunctionsNetworkDays_IntlRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) NetworkDays_Intl()(*ItemItemsItemWorkbookFunctionsNetworkDays_IntlRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNetworkDays_IntlRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Nominal provides operations to call the nominal method.
// returns a *ItemItemsItemWorkbookFunctionsNominalRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Nominal()(*ItemItemsItemWorkbookFunctionsNominalRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNominalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Norm_Dist provides operations to call the norm_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsNorm_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Norm_Dist()(*ItemItemsItemWorkbookFunctionsNorm_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNorm_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Norm_Inv provides operations to call the norm_Inv method.
// returns a *ItemItemsItemWorkbookFunctionsNorm_InvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Norm_Inv()(*ItemItemsItemWorkbookFunctionsNorm_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNorm_InvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Norm_S_Dist provides operations to call the norm_S_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsNorm_S_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Norm_S_Dist()(*ItemItemsItemWorkbookFunctionsNorm_S_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNorm_S_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Norm_S_Inv provides operations to call the norm_S_Inv method.
// returns a *ItemItemsItemWorkbookFunctionsNorm_S_InvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Norm_S_Inv()(*ItemItemsItemWorkbookFunctionsNorm_S_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNorm_S_InvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Not provides operations to call the not method.
// returns a *ItemItemsItemWorkbookFunctionsNotRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Not()(*ItemItemsItemWorkbookFunctionsNotRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNotRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Now provides operations to call the now method.
// returns a *ItemItemsItemWorkbookFunctionsNowRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Now()(*ItemItemsItemWorkbookFunctionsNowRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Nper provides operations to call the nper method.
// returns a *ItemItemsItemWorkbookFunctionsNperRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Nper()(*ItemItemsItemWorkbookFunctionsNperRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNperRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Npv provides operations to call the npv method.
// returns a *ItemItemsItemWorkbookFunctionsNpvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Npv()(*ItemItemsItemWorkbookFunctionsNpvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNpvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NumberValue provides operations to call the numberValue method.
// returns a *ItemItemsItemWorkbookFunctionsNumberValueRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) NumberValue()(*ItemItemsItemWorkbookFunctionsNumberValueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsNumberValueRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Oct2Bin provides operations to call the oct2Bin method.
// returns a *ItemItemsItemWorkbookFunctionsOct2BinRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Oct2Bin()(*ItemItemsItemWorkbookFunctionsOct2BinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOct2BinRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Oct2Dec provides operations to call the oct2Dec method.
// returns a *ItemItemsItemWorkbookFunctionsOct2DecRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Oct2Dec()(*ItemItemsItemWorkbookFunctionsOct2DecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOct2DecRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Oct2Hex provides operations to call the oct2Hex method.
// returns a *ItemItemsItemWorkbookFunctionsOct2HexRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Oct2Hex()(*ItemItemsItemWorkbookFunctionsOct2HexRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOct2HexRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Odd provides operations to call the odd method.
// returns a *ItemItemsItemWorkbookFunctionsOddRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Odd()(*ItemItemsItemWorkbookFunctionsOddRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OddFPrice provides operations to call the oddFPrice method.
// returns a *ItemItemsItemWorkbookFunctionsOddFPriceRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) OddFPrice()(*ItemItemsItemWorkbookFunctionsOddFPriceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddFPriceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OddFYield provides operations to call the oddFYield method.
// returns a *ItemItemsItemWorkbookFunctionsOddFYieldRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) OddFYield()(*ItemItemsItemWorkbookFunctionsOddFYieldRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddFYieldRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OddLPrice provides operations to call the oddLPrice method.
// returns a *ItemItemsItemWorkbookFunctionsOddLPriceRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) OddLPrice()(*ItemItemsItemWorkbookFunctionsOddLPriceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddLPriceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OddLYield provides operations to call the oddLYield method.
// returns a *ItemItemsItemWorkbookFunctionsOddLYieldRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) OddLYield()(*ItemItemsItemWorkbookFunctionsOddLYieldRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOddLYieldRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Or provides operations to call the or method.
// returns a *ItemItemsItemWorkbookFunctionsOrRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Or()(*ItemItemsItemWorkbookFunctionsOrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsOrRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the navigation property functions in drives
// returns a WorkbookFunctionsable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWorkbookFunctionsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable), nil
}
// Pduration provides operations to call the pduration method.
// returns a *ItemItemsItemWorkbookFunctionsPdurationRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Pduration()(*ItemItemsItemWorkbookFunctionsPdurationRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPdurationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Percentile_Exc provides operations to call the percentile_Exc method.
// returns a *ItemItemsItemWorkbookFunctionsPercentile_ExcRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Percentile_Exc()(*ItemItemsItemWorkbookFunctionsPercentile_ExcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPercentile_ExcRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Percentile_Inc provides operations to call the percentile_Inc method.
// returns a *ItemItemsItemWorkbookFunctionsPercentile_IncRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Percentile_Inc()(*ItemItemsItemWorkbookFunctionsPercentile_IncRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPercentile_IncRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PercentRank_Exc provides operations to call the percentRank_Exc method.
// returns a *ItemItemsItemWorkbookFunctionsPercentRank_ExcRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) PercentRank_Exc()(*ItemItemsItemWorkbookFunctionsPercentRank_ExcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPercentRank_ExcRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PercentRank_Inc provides operations to call the percentRank_Inc method.
// returns a *ItemItemsItemWorkbookFunctionsPercentRank_IncRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) PercentRank_Inc()(*ItemItemsItemWorkbookFunctionsPercentRank_IncRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPercentRank_IncRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Permut provides operations to call the permut method.
// returns a *ItemItemsItemWorkbookFunctionsPermutRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Permut()(*ItemItemsItemWorkbookFunctionsPermutRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPermutRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Permutationa provides operations to call the permutationa method.
// returns a *ItemItemsItemWorkbookFunctionsPermutationaRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Permutationa()(*ItemItemsItemWorkbookFunctionsPermutationaRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPermutationaRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Phi provides operations to call the phi method.
// returns a *ItemItemsItemWorkbookFunctionsPhiRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Phi()(*ItemItemsItemWorkbookFunctionsPhiRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPhiRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Pi provides operations to call the pi method.
// returns a *ItemItemsItemWorkbookFunctionsPiRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Pi()(*ItemItemsItemWorkbookFunctionsPiRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPiRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Pmt provides operations to call the pmt method.
// returns a *ItemItemsItemWorkbookFunctionsPmtRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Pmt()(*ItemItemsItemWorkbookFunctionsPmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPmtRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Poisson_Dist provides operations to call the poisson_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsPoisson_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Poisson_Dist()(*ItemItemsItemWorkbookFunctionsPoisson_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPoisson_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Power provides operations to call the power method.
// returns a *ItemItemsItemWorkbookFunctionsPowerRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Power()(*ItemItemsItemWorkbookFunctionsPowerRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPowerRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Ppmt provides operations to call the ppmt method.
// returns a *ItemItemsItemWorkbookFunctionsPpmtRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Ppmt()(*ItemItemsItemWorkbookFunctionsPpmtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPpmtRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Price provides operations to call the price method.
// returns a *ItemItemsItemWorkbookFunctionsPriceRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Price()(*ItemItemsItemWorkbookFunctionsPriceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPriceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PriceDisc provides operations to call the priceDisc method.
// returns a *ItemItemsItemWorkbookFunctionsPriceDiscRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) PriceDisc()(*ItemItemsItemWorkbookFunctionsPriceDiscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPriceDiscRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PriceMat provides operations to call the priceMat method.
// returns a *ItemItemsItemWorkbookFunctionsPriceMatRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) PriceMat()(*ItemItemsItemWorkbookFunctionsPriceMatRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPriceMatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Product provides operations to call the product method.
// returns a *ItemItemsItemWorkbookFunctionsProductRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Product()(*ItemItemsItemWorkbookFunctionsProductRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsProductRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Proper provides operations to call the proper method.
// returns a *ItemItemsItemWorkbookFunctionsProperRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Proper()(*ItemItemsItemWorkbookFunctionsProperRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsProperRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Pv provides operations to call the pv method.
// returns a *ItemItemsItemWorkbookFunctionsPvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Pv()(*ItemItemsItemWorkbookFunctionsPvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsPvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Quartile_Exc provides operations to call the quartile_Exc method.
// returns a *ItemItemsItemWorkbookFunctionsQuartile_ExcRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Quartile_Exc()(*ItemItemsItemWorkbookFunctionsQuartile_ExcRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsQuartile_ExcRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Quartile_Inc provides operations to call the quartile_Inc method.
// returns a *ItemItemsItemWorkbookFunctionsQuartile_IncRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Quartile_Inc()(*ItemItemsItemWorkbookFunctionsQuartile_IncRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsQuartile_IncRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Quotient provides operations to call the quotient method.
// returns a *ItemItemsItemWorkbookFunctionsQuotientRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Quotient()(*ItemItemsItemWorkbookFunctionsQuotientRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsQuotientRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Radians provides operations to call the radians method.
// returns a *ItemItemsItemWorkbookFunctionsRadiansRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Radians()(*ItemItemsItemWorkbookFunctionsRadiansRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRadiansRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Rand provides operations to call the rand method.
// returns a *ItemItemsItemWorkbookFunctionsRandRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rand()(*ItemItemsItemWorkbookFunctionsRandRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRandRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RandBetween provides operations to call the randBetween method.
// returns a *ItemItemsItemWorkbookFunctionsRandBetweenRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) RandBetween()(*ItemItemsItemWorkbookFunctionsRandBetweenRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRandBetweenRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Rank_Avg provides operations to call the rank_Avg method.
// returns a *ItemItemsItemWorkbookFunctionsRank_AvgRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rank_Avg()(*ItemItemsItemWorkbookFunctionsRank_AvgRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRank_AvgRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Rank_Eq provides operations to call the rank_Eq method.
// returns a *ItemItemsItemWorkbookFunctionsRank_EqRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rank_Eq()(*ItemItemsItemWorkbookFunctionsRank_EqRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRank_EqRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Rate provides operations to call the rate method.
// returns a *ItemItemsItemWorkbookFunctionsRateRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rate()(*ItemItemsItemWorkbookFunctionsRateRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Received provides operations to call the received method.
// returns a *ItemItemsItemWorkbookFunctionsReceivedRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Received()(*ItemItemsItemWorkbookFunctionsReceivedRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsReceivedRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Replace provides operations to call the replace method.
// returns a *ItemItemsItemWorkbookFunctionsReplaceRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Replace()(*ItemItemsItemWorkbookFunctionsReplaceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsReplaceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ReplaceB provides operations to call the replaceB method.
// returns a *ItemItemsItemWorkbookFunctionsReplaceBRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ReplaceB()(*ItemItemsItemWorkbookFunctionsReplaceBRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsReplaceBRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Rept provides operations to call the rept method.
// returns a *ItemItemsItemWorkbookFunctionsReptRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rept()(*ItemItemsItemWorkbookFunctionsReptRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsReptRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Right provides operations to call the right method.
// returns a *ItemItemsItemWorkbookFunctionsRightRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Right()(*ItemItemsItemWorkbookFunctionsRightRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRightRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Rightb provides operations to call the rightb method.
// returns a *ItemItemsItemWorkbookFunctionsRightbRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rightb()(*ItemItemsItemWorkbookFunctionsRightbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRightbRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Roman provides operations to call the roman method.
// returns a *ItemItemsItemWorkbookFunctionsRomanRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Roman()(*ItemItemsItemWorkbookFunctionsRomanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRomanRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Round provides operations to call the round method.
// returns a *ItemItemsItemWorkbookFunctionsRoundRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Round()(*ItemItemsItemWorkbookFunctionsRoundRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRoundRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RoundDown provides operations to call the roundDown method.
// returns a *ItemItemsItemWorkbookFunctionsRoundDownRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) RoundDown()(*ItemItemsItemWorkbookFunctionsRoundDownRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRoundDownRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RoundUp provides operations to call the roundUp method.
// returns a *ItemItemsItemWorkbookFunctionsRoundUpRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) RoundUp()(*ItemItemsItemWorkbookFunctionsRoundUpRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRoundUpRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Rows provides operations to call the rows method.
// returns a *ItemItemsItemWorkbookFunctionsRowsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rows()(*ItemItemsItemWorkbookFunctionsRowsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRowsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Rri provides operations to call the rri method.
// returns a *ItemItemsItemWorkbookFunctionsRriRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Rri()(*ItemItemsItemWorkbookFunctionsRriRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRriRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sec provides operations to call the sec method.
// returns a *ItemItemsItemWorkbookFunctionsSecRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sec()(*ItemItemsItemWorkbookFunctionsSecRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSecRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sech provides operations to call the sech method.
// returns a *ItemItemsItemWorkbookFunctionsSechRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sech()(*ItemItemsItemWorkbookFunctionsSechRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSechRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Second provides operations to call the second method.
// returns a *ItemItemsItemWorkbookFunctionsSecondRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Second()(*ItemItemsItemWorkbookFunctionsSecondRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSecondRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SeriesSum provides operations to call the seriesSum method.
// returns a *ItemItemsItemWorkbookFunctionsSeriesSumRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SeriesSum()(*ItemItemsItemWorkbookFunctionsSeriesSumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSeriesSumRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sheet provides operations to call the sheet method.
// returns a *ItemItemsItemWorkbookFunctionsSheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sheet()(*ItemItemsItemWorkbookFunctionsSheetRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sheets provides operations to call the sheets method.
// returns a *ItemItemsItemWorkbookFunctionsSheetsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sheets()(*ItemItemsItemWorkbookFunctionsSheetsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSheetsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sign provides operations to call the sign method.
// returns a *ItemItemsItemWorkbookFunctionsSignRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sign()(*ItemItemsItemWorkbookFunctionsSignRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSignRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sin provides operations to call the sin method.
// returns a *ItemItemsItemWorkbookFunctionsSinRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sin()(*ItemItemsItemWorkbookFunctionsSinRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSinRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sinh provides operations to call the sinh method.
// returns a *ItemItemsItemWorkbookFunctionsSinhRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sinh()(*ItemItemsItemWorkbookFunctionsSinhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSinhRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Skew provides operations to call the skew method.
// returns a *ItemItemsItemWorkbookFunctionsSkewRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Skew()(*ItemItemsItemWorkbookFunctionsSkewRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSkewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Skew_p provides operations to call the skew_p method.
// returns a *ItemItemsItemWorkbookFunctionsSkew_pRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Skew_p()(*ItemItemsItemWorkbookFunctionsSkew_pRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSkew_pRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sln provides operations to call the sln method.
// returns a *ItemItemsItemWorkbookFunctionsSlnRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sln()(*ItemItemsItemWorkbookFunctionsSlnRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSlnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Small provides operations to call the small method.
// returns a *ItemItemsItemWorkbookFunctionsSmallRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Small()(*ItemItemsItemWorkbookFunctionsSmallRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSmallRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sqrt provides operations to call the sqrt method.
// returns a *ItemItemsItemWorkbookFunctionsSqrtRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sqrt()(*ItemItemsItemWorkbookFunctionsSqrtRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSqrtRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SqrtPi provides operations to call the sqrtPi method.
// returns a *ItemItemsItemWorkbookFunctionsSqrtPiRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SqrtPi()(*ItemItemsItemWorkbookFunctionsSqrtPiRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSqrtPiRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Standardize provides operations to call the standardize method.
// returns a *ItemItemsItemWorkbookFunctionsStandardizeRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Standardize()(*ItemItemsItemWorkbookFunctionsStandardizeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStandardizeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// StDev_P provides operations to call the stDev_P method.
// returns a *ItemItemsItemWorkbookFunctionsStDev_PRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) StDev_P()(*ItemItemsItemWorkbookFunctionsStDev_PRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStDev_PRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// StDev_S provides operations to call the stDev_S method.
// returns a *ItemItemsItemWorkbookFunctionsStDev_SRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) StDev_S()(*ItemItemsItemWorkbookFunctionsStDev_SRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStDev_SRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// StDevA provides operations to call the stDevA method.
// returns a *ItemItemsItemWorkbookFunctionsStDevARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) StDevA()(*ItemItemsItemWorkbookFunctionsStDevARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStDevARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// StDevPA provides operations to call the stDevPA method.
// returns a *ItemItemsItemWorkbookFunctionsStDevPARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) StDevPA()(*ItemItemsItemWorkbookFunctionsStDevPARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsStDevPARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Substitute provides operations to call the substitute method.
// returns a *ItemItemsItemWorkbookFunctionsSubstituteRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Substitute()(*ItemItemsItemWorkbookFunctionsSubstituteRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSubstituteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Subtotal provides operations to call the subtotal method.
// returns a *ItemItemsItemWorkbookFunctionsSubtotalRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Subtotal()(*ItemItemsItemWorkbookFunctionsSubtotalRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSubtotalRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sum provides operations to call the sum method.
// returns a *ItemItemsItemWorkbookFunctionsSumRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Sum()(*ItemItemsItemWorkbookFunctionsSumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSumRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SumIf provides operations to call the sumIf method.
// returns a *ItemItemsItemWorkbookFunctionsSumIfRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SumIf()(*ItemItemsItemWorkbookFunctionsSumIfRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSumIfRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SumIfs provides operations to call the sumIfs method.
// returns a *ItemItemsItemWorkbookFunctionsSumIfsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SumIfs()(*ItemItemsItemWorkbookFunctionsSumIfsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSumIfsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SumSq provides operations to call the sumSq method.
// returns a *ItemItemsItemWorkbookFunctionsSumSqRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) SumSq()(*ItemItemsItemWorkbookFunctionsSumSqRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSumSqRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Syd provides operations to call the syd method.
// returns a *ItemItemsItemWorkbookFunctionsSydRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Syd()(*ItemItemsItemWorkbookFunctionsSydRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsSydRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// T provides operations to call the t method.
// returns a *ItemItemsItemWorkbookFunctionsTRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T()(*ItemItemsItemWorkbookFunctionsTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// T_Dist provides operations to call the t_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsT_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Dist()(*ItemItemsItemWorkbookFunctionsT_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// T_Dist_2T provides operations to call the t_Dist_2T method.
// returns a *ItemItemsItemWorkbookFunctionsT_Dist_2TRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Dist_2T()(*ItemItemsItemWorkbookFunctionsT_Dist_2TRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_Dist_2TRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// T_Dist_RT provides operations to call the t_Dist_RT method.
// returns a *ItemItemsItemWorkbookFunctionsT_Dist_RTRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Dist_RT()(*ItemItemsItemWorkbookFunctionsT_Dist_RTRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_Dist_RTRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// T_Inv provides operations to call the t_Inv method.
// returns a *ItemItemsItemWorkbookFunctionsT_InvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Inv()(*ItemItemsItemWorkbookFunctionsT_InvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_InvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// T_Inv_2T provides operations to call the t_Inv_2T method.
// returns a *ItemItemsItemWorkbookFunctionsT_Inv_2TRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) T_Inv_2T()(*ItemItemsItemWorkbookFunctionsT_Inv_2TRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsT_Inv_2TRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Tan provides operations to call the tan method.
// returns a *ItemItemsItemWorkbookFunctionsTanRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Tan()(*ItemItemsItemWorkbookFunctionsTanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTanRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Tanh provides operations to call the tanh method.
// returns a *ItemItemsItemWorkbookFunctionsTanhRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Tanh()(*ItemItemsItemWorkbookFunctionsTanhRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTanhRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TbillEq provides operations to call the tbillEq method.
// returns a *ItemItemsItemWorkbookFunctionsTbillEqRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TbillEq()(*ItemItemsItemWorkbookFunctionsTbillEqRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTbillEqRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TbillPrice provides operations to call the tbillPrice method.
// returns a *ItemItemsItemWorkbookFunctionsTbillPriceRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TbillPrice()(*ItemItemsItemWorkbookFunctionsTbillPriceRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTbillPriceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TbillYield provides operations to call the tbillYield method.
// returns a *ItemItemsItemWorkbookFunctionsTbillYieldRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TbillYield()(*ItemItemsItemWorkbookFunctionsTbillYieldRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTbillYieldRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Text provides operations to call the text method.
// returns a *ItemItemsItemWorkbookFunctionsTextRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Text()(*ItemItemsItemWorkbookFunctionsTextRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTextRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Time provides operations to call the time method.
// returns a *ItemItemsItemWorkbookFunctionsTimeRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Time()(*ItemItemsItemWorkbookFunctionsTimeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTimeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Timevalue provides operations to call the timevalue method.
// returns a *ItemItemsItemWorkbookFunctionsTimevalueRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Timevalue()(*ItemItemsItemWorkbookFunctionsTimevalueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTimevalueRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Today provides operations to call the today method.
// returns a *ItemItemsItemWorkbookFunctionsTodayRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Today()(*ItemItemsItemWorkbookFunctionsTodayRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTodayRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property functions for drives
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get functions from drives
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToPatchRequestInformation update the navigation property functions in drives
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookFunctionsable, requestConfiguration *ItemItemsItemWorkbookFunctionsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// Trim provides operations to call the trim method.
// returns a *ItemItemsItemWorkbookFunctionsTrimRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Trim()(*ItemItemsItemWorkbookFunctionsTrimRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTrimRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TrimMean provides operations to call the trimMean method.
// returns a *ItemItemsItemWorkbookFunctionsTrimMeanRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TrimMean()(*ItemItemsItemWorkbookFunctionsTrimMeanRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTrimMeanRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// True provides operations to call the true method.
// returns a *ItemItemsItemWorkbookFunctionsTrueRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) True()(*ItemItemsItemWorkbookFunctionsTrueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTrueRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Trunc provides operations to call the trunc method.
// returns a *ItemItemsItemWorkbookFunctionsTruncRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Trunc()(*ItemItemsItemWorkbookFunctionsTruncRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTruncRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TypeEscaped provides operations to call the type method.
// returns a *ItemItemsItemWorkbookFunctionsTypeRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) TypeEscaped()(*ItemItemsItemWorkbookFunctionsTypeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsTypeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Unichar provides operations to call the unichar method.
// returns a *ItemItemsItemWorkbookFunctionsUnicharRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Unichar()(*ItemItemsItemWorkbookFunctionsUnicharRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsUnicharRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Unicode provides operations to call the unicode method.
// returns a *ItemItemsItemWorkbookFunctionsUnicodeRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Unicode()(*ItemItemsItemWorkbookFunctionsUnicodeRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsUnicodeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Upper provides operations to call the upper method.
// returns a *ItemItemsItemWorkbookFunctionsUpperRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Upper()(*ItemItemsItemWorkbookFunctionsUpperRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsUpperRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Usdollar provides operations to call the usdollar method.
// returns a *ItemItemsItemWorkbookFunctionsUsdollarRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Usdollar()(*ItemItemsItemWorkbookFunctionsUsdollarRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsUsdollarRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Value provides operations to call the value method.
// returns a *ItemItemsItemWorkbookFunctionsValueRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Value()(*ItemItemsItemWorkbookFunctionsValueRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsValueRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Var_P provides operations to call the var_P method.
// returns a *ItemItemsItemWorkbookFunctionsVar_PRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Var_P()(*ItemItemsItemWorkbookFunctionsVar_PRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVar_PRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Var_S provides operations to call the var_S method.
// returns a *ItemItemsItemWorkbookFunctionsVar_SRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Var_S()(*ItemItemsItemWorkbookFunctionsVar_SRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVar_SRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// VarA provides operations to call the varA method.
// returns a *ItemItemsItemWorkbookFunctionsVarARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) VarA()(*ItemItemsItemWorkbookFunctionsVarARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVarARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// VarPA provides operations to call the varPA method.
// returns a *ItemItemsItemWorkbookFunctionsVarPARequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) VarPA()(*ItemItemsItemWorkbookFunctionsVarPARequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVarPARequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Vdb provides operations to call the vdb method.
// returns a *ItemItemsItemWorkbookFunctionsVdbRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Vdb()(*ItemItemsItemWorkbookFunctionsVdbRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVdbRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Vlookup provides operations to call the vlookup method.
// returns a *ItemItemsItemWorkbookFunctionsVlookupRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Vlookup()(*ItemItemsItemWorkbookFunctionsVlookupRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsVlookupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Weekday provides operations to call the weekday method.
// returns a *ItemItemsItemWorkbookFunctionsWeekdayRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Weekday()(*ItemItemsItemWorkbookFunctionsWeekdayRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWeekdayRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WeekNum provides operations to call the weekNum method.
// returns a *ItemItemsItemWorkbookFunctionsWeekNumRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) WeekNum()(*ItemItemsItemWorkbookFunctionsWeekNumRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWeekNumRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Weibull_Dist provides operations to call the weibull_Dist method.
// returns a *ItemItemsItemWorkbookFunctionsWeibull_DistRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Weibull_Dist()(*ItemItemsItemWorkbookFunctionsWeibull_DistRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWeibull_DistRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookFunctionsRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookFunctionsRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// WorkDay provides operations to call the workDay method.
// returns a *ItemItemsItemWorkbookFunctionsWorkDayRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) WorkDay()(*ItemItemsItemWorkbookFunctionsWorkDayRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWorkDayRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WorkDay_Intl provides operations to call the workDay_Intl method.
// returns a *ItemItemsItemWorkbookFunctionsWorkDay_IntlRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) WorkDay_Intl()(*ItemItemsItemWorkbookFunctionsWorkDay_IntlRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsWorkDay_IntlRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Xirr provides operations to call the xirr method.
// returns a *ItemItemsItemWorkbookFunctionsXirrRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Xirr()(*ItemItemsItemWorkbookFunctionsXirrRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsXirrRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Xnpv provides operations to call the xnpv method.
// returns a *ItemItemsItemWorkbookFunctionsXnpvRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Xnpv()(*ItemItemsItemWorkbookFunctionsXnpvRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsXnpvRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Xor provides operations to call the xor method.
// returns a *ItemItemsItemWorkbookFunctionsXorRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Xor()(*ItemItemsItemWorkbookFunctionsXorRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsXorRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Year provides operations to call the year method.
// returns a *ItemItemsItemWorkbookFunctionsYearRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Year()(*ItemItemsItemWorkbookFunctionsYearRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// YearFrac provides operations to call the yearFrac method.
// returns a *ItemItemsItemWorkbookFunctionsYearFracRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) YearFrac()(*ItemItemsItemWorkbookFunctionsYearFracRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYearFracRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Yield provides operations to call the yield method.
// returns a *ItemItemsItemWorkbookFunctionsYieldRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Yield()(*ItemItemsItemWorkbookFunctionsYieldRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYieldRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// YieldDisc provides operations to call the yieldDisc method.
// returns a *ItemItemsItemWorkbookFunctionsYieldDiscRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) YieldDisc()(*ItemItemsItemWorkbookFunctionsYieldDiscRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYieldDiscRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// YieldMat provides operations to call the yieldMat method.
// returns a *ItemItemsItemWorkbookFunctionsYieldMatRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) YieldMat()(*ItemItemsItemWorkbookFunctionsYieldMatRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsYieldMatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Z_Test provides operations to call the z_Test method.
// returns a *ItemItemsItemWorkbookFunctionsZ_TestRequestBuilder when successful
func (m *ItemItemsItemWorkbookFunctionsRequestBuilder) Z_Test()(*ItemItemsItemWorkbookFunctionsZ_TestRequestBuilder) {
    return NewItemItemsItemWorkbookFunctionsZ_TestRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
