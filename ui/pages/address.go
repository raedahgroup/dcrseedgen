package pages

import (
	"errors"
	"strconv"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrseedgen/helper"
	"github.com/raedahgroup/dcrseedgen/ui/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const (
	AddressPageID      = "AddressPage"
	numRowWidth        = 0.03
	addressRowWidth    = 0.4
	privateKeyRowWidth = 0.57
)

type AddressPage struct {
	theme                    *theme.Theme
	generatedAddresses       []string
	generatedPrivateKeys     []string
	numOfItemsEditorMaterial theme.Editor
	numOfItemsEditorWidget   *widget.Editor
	generateButtonMaterial   theme.Button
	generateButtonWidget     *widget.Button
	addressesLabel           theme.Label
	privateKeysLabel         theme.Label
	exportingDataLabel       theme.Label

	networkGroup         *widget.Enum
	networkRadioMaterial []theme.RadioButton

	exportIcon       theme.IconButton
	exportIconWidget *widget.Button

	message helper.Message

	isExportingData bool

	list *layout.List
	err  error
}

func NewAddressPage(th *theme.Theme) *AddressPage {
	page := &AddressPage{
		theme: th,
	}

	page.list = &layout.List{
		Axis: layout.Vertical,
	}

	networks := []string{"Testnet3", "Mainnet", "Regnet"}

	page.networkGroup = new(widget.Enum)
	page.networkGroup.SetValue(networks[0])
	page.networkRadioMaterial = make([]theme.RadioButton, len(networks))
	for i := range networks {
		page.networkRadioMaterial[i] = th.RadioButton(networks[i], networks[i])
		page.networkRadioMaterial[i].Size = unit.Dp(20)
	}

	page.numOfItemsEditorMaterial = th.Editor("How many?")
	page.numOfItemsEditorWidget = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	page.numOfItemsEditorWidget.SetText("1")

	page.generateButtonMaterial = th.Button("Generate")
	page.generateButtonWidget = new(widget.Button)

	page.addressesLabel = th.Body1("Addresses")
	page.privateKeysLabel = th.Body1("Private Keys")

	page.exportIcon = th.IconButton(theme.MustIcon(theme.NewIcon(icons.CommunicationImportExport)))
	page.exportIcon.Size = unit.Dp(30)
	page.exportIcon.Padding = unit.Dp(5)
	page.exportIconWidget = new(widget.Button)

	page.isExportingData = false
	page.exportingDataLabel = th.Caption("Exporting data...")

	return page
}

func (page *AddressPage) BeforeRender() {
	page.resetMessage()

	page.generatedAddresses = nil
	page.generatedPrivateKeys = nil

	page.numOfItemsEditorWidget.SetText("1")
}

func (page *AddressPage) resetMessage() {
	page.message.Message = ""
	page.message.Variant = ""
}

func (page *AddressPage) handleEvents(gtx *layout.Context) {
	for page.generateButtonWidget.Clicked(gtx) {
		page.resetMessage()
		page.generatePairs(page.networkGroup.Value(gtx))
	}

	for page.exportIconWidget.Clicked(gtx) {
		page.resetMessage()
		page.exportCSV()
	}
}

func (page *AddressPage) exportCSV() {
	now := time.Now()
	filename, err := time.Parse(time.RFC3339, now.Format(time.RFC3339))
	if err != nil {
		panic(err)
	}

	page.isExportingData = true

	// prepare data
	data := make([][]string, len(page.generatedAddresses))
	for index := range page.generatedAddresses {
		data[index] = []string{page.generatedAddresses[index], page.generatedPrivateKeys[index]}
	}

	// show exporting message
	exportPath, err := helper.CreateCSV(filename.String(), data)
	if err != nil {
		page.message.Message = "error exporting data: " + err.Error()
		page.message.Variant = "error"
	} else {
		page.message.Message = "Successfully exported csv to " + exportPath
		page.message.Variant = "success"
	}
	page.isExportingData = false
}

func (page *AddressPage) generatePairs(network string) {
	numberOfItemsToGenerateStr := page.numOfItemsEditorWidget.Text()
	if numberOfItemsToGenerateStr == "" {
		page.err = errors.New("Please type in the required number of pairs")
		return
	}

	numberOfItemsToGenerate, err := strconv.Atoi(numberOfItemsToGenerateStr)
	if err != nil {
		page.err = errors.New("Invalid number")
		return
	}

	page.generatedAddresses = make([]string, numberOfItemsToGenerate)
	page.generatedPrivateKeys = make([]string, numberOfItemsToGenerate)

	for i := 0; i < numberOfItemsToGenerate; i++ {
		privateKey, address, err := helper.GenerateAddressAndPrivateKey(network)
		if err != nil {
			page.err = err
			return
		}

		page.generatedAddresses[i] = address
		page.generatedPrivateKeys[i] = privateKey
	}
}

func (page *AddressPage) Render(gtx *layout.Context) {
	page.handleEvents(gtx)

	layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				layout.Inset{Bottom: unit.Dp(15)}.Layout(gtx, func() {
					page.renderFormSection(gtx)
				})
			}),
			layout.Rigid(func() {
				if page.message.Message != "" {
					if page.message.Variant == "error" {
						page.theme.ErrorAlert(gtx, page.message.Message)
					} else {
						page.theme.SuccessAlert(gtx, page.message.Message)
					}
				}
			}),
			layout.Rigid(func() {
				if len(page.generatedAddresses) > 0 {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func() {
								page.renderHeader(gtx)
							})
						}),
						layout.Flexed(0.85, func() {
							page.list.Layout(gtx, len(page.generatedAddresses), func(i int) {
								page.renderRow(gtx, i)
							})
						}),
					)
				}
			}),
			layout.Rigid(func() {
				if len(page.generatedAddresses) > 0 {
					layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
						layout.Rigid(func() {
							page.exportIcon.Layout(gtx, page.exportIconWidget)
						}),
						layout.Rigid(func() {
							layout.Inset{Top: unit.Dp(3)}.Layout(gtx, func() {
								layout.Center.Layout(gtx, func() {
									page.theme.Caption("Export").Layout(gtx)
								})
							})
						}),
					)
				}

			}),
		)
	})

	if page.isExportingData {
		page.renderExportingModal(gtx)
	}
}

func (page *AddressPage) renderHeader(gtx *layout.Context) {
	txt := page.theme.Label(unit.Dp(16), "#")
	txt.Color = page.theme.Color.Hint

	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(numRowWidth, func() {
			txt.Layout(gtx)
		}),
		layout.Flexed(addressRowWidth, func() {
			txt.Text = "Address"
			txt.Layout(gtx)
		}),
		layout.Flexed(privateKeyRowWidth, func() {
			txt.Text = "Private Key"
			txt.Layout(gtx)
		}),
	)
}

func (page *AddressPage) renderRow(gtx *layout.Context, index int) {
	layout.Inset{Bottom: unit.Dp(15)}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(numRowWidth, func() {
				page.theme.Caption(strconv.Itoa(index + 1)).Layout(gtx)
			}),
			layout.Flexed(addressRowWidth, func() {
				page.theme.Caption(page.generatedAddresses[index]).Layout(gtx)
			}),
			layout.Flexed(privateKeyRowWidth, func() {
				page.theme.Caption(page.generatedPrivateKeys[index]).Layout(gtx)
			}),
		)
	})
}

func (page *AddressPage) renderGeneratedPairs(gtx *layout.Context) {
	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, len(page.generatedAddresses), func(i int) {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func() {
				page.theme.Body2(page.generatedAddresses[i]).Layout(gtx)
			}),
			layout.Rigid(func() {
				page.theme.Body2(page.generatedPrivateKeys[i]).Layout(gtx)
			}),
		)
	})
}

func (page *AddressPage) renderFormSection(gtx *layout.Context) {
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func() {
			(&layout.List{Axis: layout.Horizontal}).
				Layout(gtx, len(page.networkRadioMaterial), func(index int) {
					layout.Inset{Right: unit.Dp(5), Top: unit.Dp(15)}.Layout(gtx, func() {
						page.networkRadioMaterial[index].Layout(gtx, page.networkGroup)
					})
				})
		}),
		layout.Rigid(func() {
			page.drawDivider(gtx)
		}),
		layout.Rigid(func() {
			gtx.Constraints.Width.Max = 100
			page.numOfItemsEditorMaterial.Layout(gtx, page.numOfItemsEditorWidget)
		}),
		layout.Rigid(func() {
			page.drawDivider(gtx)
		}),
		layout.Rigid(func() {
			page.generateButtonMaterial.Layout(gtx, page.generateButtonWidget)
		}),
	)
}

func (page *AddressPage) drawDivider(gtx *layout.Context) {
	layout.Inset{
		Left:  unit.Dp(25),
		Right: unit.Dp(25),
		Top:   unit.Dp(0)}.Layout(gtx, func() {
		rect := f32.Rectangle{Max: f32.Point{X: 1, Y: 45}}
		op.TransformOp{}.Offset(f32.Point{X: 0, Y: 0}).Add(gtx.Ops)
		paint.ColorOp{Color: page.theme.Color.Hint}.Add(gtx.Ops)
		paint.PaintOp{Rect: rect}.Add(gtx.Ops)
	})
}

func (page *AddressPage) renderExportingModal(gtx *layout.Context) {
	overlayColor := page.theme.Color.Black
	overlayColor.A = 200

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			theme.FillMax(gtx, overlayColor)
		}),
		layout.Stacked(func() {
			hv := float32((gtx.Constraints.Height.Max / 2) - 30)

			layout.Inset{
				Top:    unit.Dp(hv),
				Bottom: unit.Dp(hv),
				Left:   unit.Dp(80),
				Right:  unit.Dp(80),
			}.Layout(gtx, func() {
				theme.FillMax(gtx, page.theme.Color.Surface)

				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				gtx.Constraints.Height.Min = gtx.Constraints.Height.Max

				layout.Center.Layout(gtx, func() {
					page.exportingDataLabel.Layout(gtx)
				})
			})
		}),
	)
}
