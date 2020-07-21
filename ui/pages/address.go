package pages

import (
	"errors"
	"image"
	"strconv"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

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
	generateButtonWidget     *widget.Clickable
	addressesLabel           material.LabelStyle
	privateKeysLabel         material.LabelStyle
	exportingDataLabel       material.LabelStyle

	networkGroup         *widget.Enum
	networkRadioMaterial []theme.RadioButton

	exportIcon       theme.IconButton
	exportIconWidget *widget.Clickable

	message helper.Message

	isExportingData bool

	list        *layout.List
	addressList *layout.List
	err         error
}

func NewAddressPage(th *theme.Theme) *AddressPage {
	page := &AddressPage{
		theme: th,
	}

	page.list = &layout.List{
		Axis: layout.Vertical,
	}

	page.addressList = &layout.List{
		Axis: layout.Vertical,
	}

	networks := []string{"Testnet3", "Mainnet", "Regnet"}

	page.networkGroup = new(widget.Enum)
	page.networkGroup.Value = networks[0]
	page.networkRadioMaterial = make([]theme.RadioButton, len(networks))
	for i := range networks {
		page.networkRadioMaterial[i] = th.RadioButton(networks[i], networks[i], page.networkGroup)
		page.networkRadioMaterial[i].Size = unit.Dp(20)
	}

	page.numOfItemsEditorWidget = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	page.numOfItemsEditorMaterial = th.Editor("How many?", page.numOfItemsEditorWidget)
	page.numOfItemsEditorWidget.SetText("1")

	page.generateButtonWidget = new(widget.Clickable)
	page.generateButtonMaterial = th.Button("Generate", page.generateButtonWidget)

	page.addressesLabel = th.Body1("Addresses")
	page.privateKeysLabel = th.Body1("Private Keys")

	page.exportIconWidget = new(widget.Clickable)
	page.exportIcon = th.IconButton(theme.MustIcon(theme.NewIcon(icons.CommunicationImportExport)), page.exportIconWidget)
	page.exportIcon.Size = unit.Dp(30)
	page.exportIcon.Padding = unit.Dp(5)
	page.exportIcon.Color = th.Color.Surface

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

func (page *AddressPage) handleEvents() {
	for page.generateButtonWidget.Clicked() {
		page.resetMessage()
		page.generatePairs(page.networkGroup.Value)
	}

	for page.exportIconWidget.Clicked() {
		page.resetMessage()
		page.exportCSV()
	}

	if page.message.Message != "" && page.message.Variant == "success" {
		time.AfterFunc(time.Second*5, func() {
			page.message.Message = ""
		})
	}
}

func (page *AddressPage) exportCSV() {
	page.isExportingData = true

	// prepare data
	data := make([][]string, len(page.generatedAddresses))
	for index := range page.generatedAddresses {
		data[index] = []string{page.generatedAddresses[index], page.generatedPrivateKeys[index]}
	}

	// show exporting message
	exportPath, err := helper.CreateCSV(data)
	if err != nil {
		page.message.Message = "error exporting data: " + err.Error()
		page.message.Variant = "error"
	} else {
		page.message.Message = "Exported data to " + exportPath
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
		page.err = errors.New("Please specify a valid number to generate")
		return
	}

	page.err = nil

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

func (page *AddressPage) Render(gtx layout.Context) layout.Dimensions {
	page.handleEvents()
	maxHeight := gtx.Constraints.Max.Y

	w := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return page.renderFormSection(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			if page.err != nil {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return page.theme.ErrorAlert(gtx, page.err.Error())
			}

			return layout.Dimensions{}
		},
		func(gtx layout.Context) layout.Dimensions {
			if page.message.Message != "" {
				if page.message.Variant == "error" {
					return page.theme.ErrorAlert(gtx, page.message.Message)
				}
				return page.theme.SuccessAlert(gtx, page.message.Message)
			}
			return layout.Dimensions{}
		},
		func(gtx layout.Context) layout.Dimensions {

			if len(page.generatedAddresses) > 0 {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return page.renderHeader(gtx)
						})
					}),
					layout.Flexed(0.8, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = int(float32(0.67) * float32(maxHeight))
						return page.addressList.Layout(gtx, len(page.generatedAddresses), func(gtx layout.Context, i int) layout.Dimensions {
							return page.renderRow(gtx, i)
						})
					}),
				)
			}
			return layout.Dimensions{}
		},
		func(gtx layout.Context) layout.Dimensions {
			if len(page.generatedAddresses) > 0 {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return page.exportIcon.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(3)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return page.theme.Caption("Export").Layout(gtx)
							})
						})
					}),
				)
			}
			return layout.Dimensions{}
		},
	}

	if page.isExportingData {
		page.renderExportingModal(gtx)
	}

	return page.list.Layout(gtx, len(w), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, w[i])
	})
}

func (page *AddressPage) renderHeader(gtx layout.Context) layout.Dimensions {
	txt := page.theme.Label(unit.Dp(16), "#")
	txt.Color = page.theme.Color.Hint

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(numRowWidth, func(gtx layout.Context) layout.Dimensions {
			return txt.Layout(gtx)
		}),
		layout.Flexed(addressRowWidth, func(gtx layout.Context) layout.Dimensions {
			txt.Text = "Address"
			return txt.Layout(gtx)
		}),
		layout.Flexed(privateKeyRowWidth, func(gtx layout.Context) layout.Dimensions {
			txt.Text = "Private Key"
			return txt.Layout(gtx)
		}),
	)
}

func (page *AddressPage) renderRow(gtx layout.Context, index int) layout.Dimensions {
	return layout.Inset{Bottom: unit.Dp(15)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(numRowWidth, func(gtx layout.Context) layout.Dimensions {
				return page.theme.Caption(strconv.Itoa(index + 1)).Layout(gtx)
			}),
			layout.Flexed(addressRowWidth, func(gtx layout.Context) layout.Dimensions {
				return page.theme.Caption(page.generatedAddresses[index]).Layout(gtx)
			}),
			layout.Flexed(privateKeyRowWidth, func(gtx layout.Context) layout.Dimensions {
				return page.theme.Caption(page.generatedPrivateKeys[index]).Layout(gtx)
			}),
		)
	})
}

func (page *AddressPage) renderGeneratedPairs(gtx layout.Context) layout.Dimensions {
	list := layout.List{Axis: layout.Vertical}
	return list.Layout(gtx, len(page.generatedAddresses), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return page.theme.Body2(page.generatedAddresses[i]).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return page.theme.Body2(page.generatedPrivateKeys[i]).Layout(gtx)
			}),
		)
	})
}

func (page *AddressPage) renderFormSection(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			list := layout.List{Axis: layout.Horizontal}
			return list.Layout(gtx, len(page.networkRadioMaterial), func(gtx layout.Context, index int) layout.Dimensions {
				return layout.Inset{Right: unit.Dp(5), Top: unit.Dp(15)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return page.networkRadioMaterial[index].Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return page.drawDivider(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 100
			return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return page.numOfItemsEditorMaterial.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return page.drawDivider(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return page.generateButtonMaterial.Layout(gtx)
		}),
	)
}

func (page *AddressPage) drawDivider(gtx layout.Context) layout.Dimensions {
	x, y := 1, 45

	return layout.Inset{
		Left:  unit.Dp(25),
		Right: unit.Dp(25),
		Top:   unit.Dp(0),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		rect := f32.Rectangle{Max: f32.Point{X: float32(x), Y: float32(y)}}
		op.TransformOp{}.Offset(f32.Point{X: 0, Y: 0}).Add(gtx.Ops)
		paint.ColorOp{Color: page.theme.Color.Hint}.Add(gtx.Ops)
		paint.PaintOp{Rect: rect}.Add(gtx.Ops)

		return layout.Dimensions{
			Size: image.Point{X: x, Y: y},
		}
	})
}

func (page *AddressPage) renderExportingModal(gtx layout.Context) layout.Dimensions {
	overlayColor := page.theme.Color.Black
	overlayColor.A = 200

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return theme.FillMax(gtx, overlayColor)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			hv := float32((gtx.Constraints.Max.Y / 2) - 30)

			return layout.Inset{
				Top:    unit.Dp(hv),
				Bottom: unit.Dp(hv),
				Left:   unit.Dp(80),
				Right:  unit.Dp(80),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				theme.FillMax(gtx, page.theme.Color.Surface)

				gtx.Constraints.Min = gtx.Constraints.Max

				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return page.exportingDataLabel.Layout(gtx)
				})
			})
		}),
	)
}
