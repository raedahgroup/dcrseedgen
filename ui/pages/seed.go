package pages

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

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
		seedWordsHeaderLabel        theme.Label
		seedHexHeaderLabel          theme.Label
		seedVerificationHeaderLabel theme.Label

		isShowingVerificationPage bool

		verifyButtonMaterial theme.Button
		verifyButtonWidget   *widget.Button

		backVerificationButtonMaterial theme.Button
		backVerificationButtonWidget   *widget.Button

		doVerifyButtonMaterial theme.Button
		doVerifyButtonWidget   *widget.Button

		generateButtonMaterial theme.Button
		generateButtonWidget   *widget.Button

		copyIconMaterial theme.IconButton
		copyIconWidget   *widget.Button

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
		theme: th,
	}

	page.list = &layout.List{
		Axis: layout.Vertical,
	}

	page.verifyMessage = helper.Message{}

	page.err = errors.New("dddd")
	page.hasCopiedHexSeed = false
	page.isShowingVerificationPage = false
	page.seedWordsHeaderLabel = th.H5("Seed Words")
	page.seedHexHeaderLabel = th.H5("Seed Hex")
	page.verifyButtonMaterial = th.Button("Verify")
	page.generateButtonMaterial = th.Button("Regenerate")
	page.seedVerificationHeaderLabel = th.H5("Verify Seed Words")
	page.verifyButtonWidget = new(widget.Button)
	page.generateButtonWidget = new(widget.Button)
	page.copyIconWidget = new(widget.Button)
	page.copyIconMaterial = th.IconButton(theme.MustIcon(theme.NewIcon(icons.ContentContentCopy)))
	page.copyIconMaterial.Background = th.Color.Background
	page.copyIconMaterial.Color = th.Color.Text
	page.copyIconMaterial.Size = unit.Dp(25)
	page.copyIconMaterial.Padding = unit.Dp(5)

	page.doVerifyButtonMaterial = th.Button("Verify")
	page.doVerifyButtonWidget = new(widget.Button)

	page.backVerificationButtonMaterial = th.Button("Back")
	page.backVerificationButtonWidget = new(widget.Button)

	page.generate()
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

func (page *SeedPage) handleEvents(gtx *layout.Context) {
	if page.hasCopiedHexSeed {
		time.AfterFunc(3*time.Second, func() {
			page.hasCopiedHexSeed = false
		})
	}

	for page.generateButtonWidget.Clicked(gtx) {
		page.generate()
	}

	for page.copyIconWidget.Clicked(gtx) {
		clipboard.WriteAll(page.seed.seedStr)
		page.hasCopiedHexSeed = true
	}

	for page.verifyButtonWidget.Clicked(gtx) {
		page.isShowingVerificationPage = true
	}
}

func (page *SeedPage) Render(gtx *layout.Context) {
	page.handleEvents(gtx)

	if page.isShowingVerificationPage {
		page.renderSeedVerificationPage(gtx)
	} else {
		page.renderSeedGenerationPage(gtx)
	}
}

func (page *SeedPage) renderSeedGenerationPage(gtx *layout.Context) {
	if page.err != nil {
		page.theme.ErrorAlert(gtx, page.err.Error())
		return
	}

	w := []func(){
		func() {
			page.seedWordsHeaderLabel.Layout(gtx)
		},
		func() {
			layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func() {
				page.renderWordColumns(gtx)
			})
		},
		func() {
			layout.Inset{Top: unit.Dp(30)}.Layout(gtx, func() {
				page.seedHexHeaderLabel.Layout(gtx)
			})
		},
		func() {
			layout.Inset{Top: unit.Dp(5)}.Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						page.theme.Body1(page.seed.seedStr).Layout(gtx)
					}),
					layout.Rigid(func() {
						layout.Inset{Top: unit.Dp(-3), Left: unit.Dp(3)}.Layout(gtx, func() {
							page.copyIconMaterial.Layout(gtx, page.copyIconWidget)
						})
						if page.hasCopiedHexSeed {
							layout.Inset{Top: unit.Dp(25)}.Layout(gtx, func() {
								page.theme.Caption("copied").Layout(gtx)
							})
						}
					}),
				)
			})
		},
		func() {
			insetTop := float32(35)
			if page.hasCopiedHexSeed {
				insetTop = 15
			}

			layout.Inset{Top: unit.Dp(insetTop)}.Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func() {
						page.generateButtonMaterial.Layout(gtx, page.generateButtonWidget)
					}),
					layout.Rigid(func() {
						page.verifyButtonMaterial.Layout(gtx, page.verifyButtonWidget)
					}),
				)
			})
		},
	}

	page.list.Layout(gtx, len(w), func(i int) {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, w[i])
	})
}

func (page *SeedPage) renderWordColumns(gtx *layout.Context) {
	currentItem := 1

	(&layout.List{Axis: layout.Horizontal}).Layout(gtx, len(page.seed.columns), func(i int) {
		layout.Inset{
			Right: unit.Dp(30),
		}.Layout(gtx, func() {
			(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(page.seed.columns[i].words), func(j int) {
				layout.Inset{
					Bottom: unit.Dp(10),
				}.Layout(gtx, func() {
					page.theme.Body1(strconv.Itoa(currentItem) + ". " + page.seed.columns[i].words[j]).Layout(gtx)
				})
				currentItem++
			})
		})
	})
}

func (page *SeedPage) renderSeedVerificationPage(gtx *layout.Context) {
	page.handleVerificationEvents(gtx)

	w := []func(){
		func() {
			page.seedVerificationHeaderLabel.Layout(gtx)
		},
		func() {
			if page.verifyMessage.Message != "" {
				if page.verifyMessage.Variant == "success" {
					page.theme.SuccessAlert(gtx, page.verifyMessage.Message)
				} else {
					page.theme.ErrorAlert(gtx, page.verifyMessage.Message)
				}
			}
		},
		func() {
			layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func() {
				page.renderInputColumns(gtx)
			})
		},
		func() {
			layout.Inset{Top: unit.Dp(30)}.Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func() {
						page.backVerificationButtonMaterial.Layout(gtx, page.backVerificationButtonWidget)
					}),
					layout.Rigid(func() {
						page.doVerifyButtonMaterial.Layout(gtx, page.doVerifyButtonWidget)
					}),
				)
			})
		},
	}

	page.list.Layout(gtx, len(w), func(i int) {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, w[i])
	})
}

func (page *SeedPage) handleVerificationEvents(gtx *layout.Context) {
	for page.backVerificationButtonWidget.Clicked(gtx) {
		page.isShowingVerificationPage = false
		page.resetVerificationPage()
	}

	for page.doVerifyButtonWidget.Clicked(gtx) {
		page.doVerification()
	}
}

func (page *SeedPage) renderInputColumns(gtx *layout.Context) {
	currentItem := 1
	maxWidth := gtx.Constraints.Width.Max

	(&layout.List{Axis: layout.Horizontal}).Layout(gtx, len(page.seed.columns), func(i int) {
		layout.Inset{
			Right: unit.Dp(30),
		}.Layout(gtx, func() {
			(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(page.seed.columns[i].editors), func(j int) {
				layout.Inset{
					Bottom: unit.Dp(10),
				}.Layout(gtx, func() {
					gtx.Constraints.Width.Max = maxWidth / (numberOfColumns + 1)
					layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func() {
							page.theme.Body1(strconv.Itoa(currentItem)).Layout(gtx)
						}),
						layout.Rigid(func() {
							ed := page.theme.Editor("")
							ed.IsTitleLabel = false
							ed.Layout(gtx, page.seed.columns[i].editors[j])
						}),
					)
				})
				currentItem++
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
