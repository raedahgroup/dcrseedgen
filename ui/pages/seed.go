package pages

import (
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/atotto/clipboard"
	"github.com/raedahgroup/dcrseedgen/helper"
	"github.com/raedahgroup/dcrseedgen/ui/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type (
	column struct {
		words   []string
		editors []*widget.Editor
	}

	seed struct {
		seedStr string
		columns []column
	}

	SeedPage struct {
		theme       *theme.Theme
		currentPage string
		seed        *seed
		err         error

		hasCopiedHexSeed bool

		list                        *layout.List
		seedWordsHeaderLabel        material.LabelStyle
		seedHexHeaderLabel          material.LabelStyle
		seedVerificationHeaderLabel material.LabelStyle
		copiedLabel                 material.LabelStyle

		isShowingVerificationPage bool

		verifyButtonMaterial theme.Button
		verifyButtonWidget   *widget.Clickable

		backVerificationButtonMaterial theme.Button
		backVerificationButtonWidget   *widget.Clickable

		doVerifyButtonMaterial theme.Button
		doVerifyButtonWidget   *widget.Clickable

		generateButtonMaterial theme.Button
		generateButtonWidget   *widget.Clickable

		copyIconMaterial theme.IconButton
		copyIconWidget   *widget.Clickable

		verifyMessage helper.Message
	}
)

const (
	SeedPageID = "SeedPage"

	seedSize        = 32 // affects number of words // noOfWords = (seedSize+1)
	numberOfColumns = 5
	numberOfRows    = 7
)

func NewSeedPage(th *theme.Theme) *SeedPage {
	page := &SeedPage{
		theme:         th,
		verifyMessage: helper.Message{},
	}

	page.list = &layout.List{
		Axis: layout.Vertical,
	}

	page.hasCopiedHexSeed = false
	page.isShowingVerificationPage = false

	page.seedWordsHeaderLabel = th.H5("Seed Words")
	page.seedHexHeaderLabel = th.H5("Seed Hex")
	page.copiedLabel = th.Caption("copied")
	page.seedVerificationHeaderLabel = th.H5("Verify Seed Words")

	page.verifyButtonWidget = new(widget.Clickable)
	page.verifyButtonMaterial = th.Button("Verify", page.verifyButtonWidget)

	page.generateButtonWidget = new(widget.Clickable)
	page.generateButtonMaterial = th.Button("Regenerate", page.generateButtonWidget)

	page.copyIconWidget = new(widget.Clickable)
	page.copyIconMaterial = th.IconButton(theme.MustIcon(theme.NewIcon(icons.ContentContentCopy)), page.copyIconWidget)
	page.copyIconMaterial.Background = th.Color.Background
	page.copyIconMaterial.Color = th.Color.Text
	page.copyIconMaterial.Size = unit.Dp(25)
	page.copyIconMaterial.Padding = unit.Dp(5)

	page.doVerifyButtonWidget = new(widget.Clickable)
	page.doVerifyButtonMaterial = th.Button("Verify", page.doVerifyButtonWidget)

	page.backVerificationButtonWidget = new(widget.Clickable)
	page.backVerificationButtonMaterial = th.DangerButton("Back", page.backVerificationButtonWidget)

	page.copiedLabel.Color = th.Color.Success

	return page
}

func (page *SeedPage) generate() {
	words, seedStr, err := helper.GenerateMnemonicSeed(seedSize)
	if err != nil {
		page.err = err
		return
	}

	page.seed = &seed{
		seedStr: seedStr,
		columns: make([]column, numberOfColumns),
	}

	wordSlice := strings.Split(words, " ")
	currentColumn := 0

	for index, word := range wordSlice {
		page.seed.columns[currentColumn].words = append(page.seed.columns[currentColumn].words, word)
		editor := &widget.Editor{
			SingleLine: true,
			Submit:     true,
		}
		page.seed.columns[currentColumn].editors = append(page.seed.columns[currentColumn].editors, editor)

		if index > 0 && (index+1)%numberOfRows == 0 {
			currentColumn++
		}
	}
}

func (page *SeedPage) BeforeRender() {
	page.resetSeedGenerationPage()
	page.resetVerificationPage()
}

func (page *SeedPage) handleEvents(gtx layout.Context) {
	if page.hasCopiedHexSeed {
		time.AfterFunc(3*time.Second, func() {
			page.hasCopiedHexSeed = false
		})
	}

	for page.generateButtonWidget.Clicked() {
		page.generate()
	}

	for page.copyIconWidget.Clicked() {
		clipboard.WriteAll(page.seed.seedStr)
		page.hasCopiedHexSeed = true
	}

	for page.verifyButtonWidget.Clicked() {
		page.isShowingVerificationPage = true
	}
}

func (page *SeedPage) Render(gtx layout.Context) layout.Dimensions {
	page.handleEvents(gtx)

	if page.isShowingVerificationPage {
		return page.renderSeedVerificationPage(gtx)
	}
	return page.renderSeedGenerationPage(gtx)
}

func (page *SeedPage) renderSeedGenerationPage(gtx layout.Context) layout.Dimensions {
	if page.err != nil {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return page.theme.ErrorAlert(gtx, page.err.Error())
	}

	w := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return page.seedWordsHeaderLabel.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return page.renderWordColumns(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return page.seedHexHeaderLabel.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return page.theme.Body1(page.seed.seedStr).Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Top: unit.Dp(0), Left: unit.Dp(7)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return page.copyIconMaterial.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if page.hasCopiedHexSeed {
								return layout.Inset{Top: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return page.copiedLabel.Layout(gtx)
								})
							}
							return layout.Dimensions{}
						}),
					)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			insetTop := float32(35)
			if page.hasCopiedHexSeed {
				insetTop = 20
			}

			return layout.Inset{Top: unit.Dp(insetTop)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return page.generateButtonMaterial.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return page.verifyButtonMaterial.Layout(gtx)
					}),
				)
			})
		},
	}

	return page.list.Layout(gtx, len(w), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, w[i])
	})
}

func (page *SeedPage) renderWordColumns(gtx layout.Context) layout.Dimensions {
	colWidth := gtx.Constraints.Max.X / len(page.seed.columns)

	currentItem := 1
	columnList := layout.List{Axis: layout.Horizontal}
	return columnList.Layout(gtx, len(page.seed.columns), func(gtx layout.Context, i int) layout.Dimensions {
		wordList := layout.List{Axis: layout.Vertical}
		return wordList.Layout(gtx, len(page.seed.columns[i].words), func(gtx layout.Context, j int) layout.Dimensions {
			w := layout.Inset{
				Bottom: unit.Dp(10),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = colWidth
				return page.theme.Body1(strconv.Itoa(currentItem) + ". " + page.seed.columns[i].words[j]).Layout(gtx)
			})

			currentItem++
			return w
		})
	})
}

func (page *SeedPage) renderSeedVerificationPage(gtx layout.Context) layout.Dimensions {
	page.handleVerificationEvents()

	w := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return page.seedVerificationHeaderLabel.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			if page.verifyMessage.Message != "" {
				if page.verifyMessage.Variant == "success" {
					return page.theme.SuccessAlert(gtx, page.verifyMessage.Message)
				} else {
					return page.theme.ErrorAlert(gtx, page.verifyMessage.Message)
				}
			}

			return layout.Dimensions{}
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return page.renderInputColumns(gtx)
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(30)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return page.backVerificationButtonMaterial.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return page.doVerifyButtonMaterial.Layout(gtx)
					}),
				)
			})
		},
	}

	return page.list.Layout(gtx, len(w), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, w[i])
	})
}

func (page *SeedPage) handleVerificationEvents() {
	for page.backVerificationButtonWidget.Clicked() {
		page.isShowingVerificationPage = false
		page.resetVerificationPage()
	}

	for page.doVerifyButtonWidget.Clicked() {
		page.doVerification()
	}
}

func (page *SeedPage) renderInputColumns(gtx layout.Context) layout.Dimensions {
	currentItem := 1
	maxWidth := gtx.Constraints.Max.X

	return (&layout.List{Axis: layout.Horizontal}).Layout(gtx, len(page.seed.columns), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.Inset{
			Right: unit.Dp(30),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(page.seed.columns[i].editors), func(gtx layout.Context, j int) layout.Dimensions {
				dims := layout.Inset{
					Bottom: unit.Dp(10),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = maxWidth / (numberOfColumns + 1)
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return page.theme.Body1(strconv.Itoa(currentItem)).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Inset{Left: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return page.theme.Editor("", page.seed.columns[i].editors[j]).Layout(gtx)
							})
						}),
					)
				})
				currentItem++
				return dims
			})
		})
	})
}

func (page *SeedPage) doVerification() {
	for range page.seed.columns {
		for columnIndex := range page.seed.columns {
			for itemIndex := range page.seed.columns[columnIndex].words {
				if page.seed.columns[columnIndex].words[itemIndex] != page.seed.columns[columnIndex].editors[itemIndex].Text() {
					page.verifyMessage.Message = "Invalid verification words. Please check the words and try again"
					page.verifyMessage.Variant = "error"
					return
				}
			}
		}
	}
	page.verifyMessage.Message = "Verification successfull"
	page.verifyMessage.Variant = "success"
}

func (page *SeedPage) resetVerificationPage() {
	page.verifyMessage.Message = ""
	page.verifyMessage.Variant = ""

	for range page.seed.columns {
		for columnIndex := range page.seed.columns {
			for itemIndex := range page.seed.columns[columnIndex].words {
				page.seed.columns[columnIndex].editors[itemIndex].SetText("")
			}
		}
	}
}

func (page *SeedPage) resetSeedGenerationPage() {
	page.err = nil

	if page.isShowingVerificationPage {
		return
	}

	page.seed = nil
	page.generate()
}
