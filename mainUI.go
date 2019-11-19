package main

import (
	"log"
	"github.com/lxn/walk"
	."github.com/lxn/walk/declarative"
)

type DVWAfker struct {
	window         *walk.MainWindow
	urlLabel	*walk.Label
	dvwaEdit		   *walk.LineEdit
	debugInfo       *walk.Label
	confirmButton		*walk.PushButton
}

func createDVWA() (*DVWAfker, error) {
	var df DVWAfker
	def := MainWindow{
		AssignTo: &df.window,
		Title:    "dvwa crack",
		MinSize:  Size{Width: 320, Height: 300},
		Size:     Size{Width: 640, Height: 600},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout :Grid{Columns:2,Spacing:10},
				Children: []Widget{
					VSplitter{
						Children:	[]Widget{
							Label{
								AssignTo:	&df.urlLabel,
								Text:	"DVWA URL:",
								Visible: true,
							},
							LineEdit{
								MinSize:	Size{Width: 120,Height: 0},
								AssignTo:	&df.dvwaEdit,
							},
							PushButton{
								Text:    "fkit!",
								MinSize: Size{Width: 120,Height: 50},
								OnClicked: func() {
									if df.dvwaEdit.ReadOnly(){
										var tmp walk.Form
										walk.MsgBox(tmp, "Warning", "Already init!", walk.MsgBoxIconInformation)
										return
									}
									if df.dvwaEdit.Text() == ""{
										var tmp walk.Form
										walk.MsgBox(tmp, "Warning", "Can not be empty!", walk.MsgBoxIconInformation)
										return
									}else{
										cookie := InitDvwaUrl(df.dvwaEdit.Text())
										if cookie == ""{
											var tmp walk.Form
											walk.MsgBox(tmp, "Error", "Can not get cookie!", walk.MsgBoxIconInformation)
											return
										}else{
											var tmp walk.Form
											walk.MsgBox(tmp, "Info", "Success init!", walk.MsgBoxIconInformation)
											df.urlLabel.SetText("PHPSSID:")
											df.dvwaEdit.SetText(cookie)
											df.dvwaEdit.SetReadOnly(true)
											return
										}
									}
								},
							},
						},
					},
				},
			},
		},
	}
	err := def.Create()
	if err != nil {
		return nil, err
	}
	return &df, nil
}


func doRunGUI() {
	fm, err := createDVWA()
	if err != nil {
		log.Fatalf("failed with '%s'\n", err)
	}
	_ = fm.window.Run()
}

func main(){
	doRunGUI()
}
