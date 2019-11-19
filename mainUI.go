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
	itemInfo		*walk.Label
	debugInfo       *walk.Label
	confirmButton		*walk.PushButton

	urlDVWA	string
	cookieDVWA	string
	confirmRes	string
}

type subWindow struct{
		//subWindow
		editFile	 *walk.LineEdit
		fileButton		*walk.PushButton
		attackButton		*walk.PushButton
		filePath 		string
		outPut	 *walk.TextEdit
		*walk.MainWindow
}

func (sw *subWindow) pbClicked()error {

    dlg := new(walk.FileDialog)
    dlg.FilePath = sw.filePath
    dlg.Title = "Select File"
	dlg.Filter = "Txt files (*.txt)|*.txt|All files (*.*)|*.*"
	sw.editFile.SetText("")
	if ok,err:= dlg.ShowOpen(sw); err!=nil {
        return err
    } else if!ok {
        return nil
    }
	sw.editFile.SetText(dlg.FilePath)
	sw.filePath = dlg.FilePath
	return nil
}


func createDVWA() (*DVWAfker, error) {
	df := &DVWAfker{}
	sw := &subWindow{}
	subdef := MainWindow{
		AssignTo: &sw.MainWindow,
		Title:    "dvwa crack",
		MinSize:  Size{Width: 200, Height: 250},
		Size:     Size{Width: 400, Height: 500},
		Layout:	VBox{},
		Children:	[]Widget{
				Composite{
					Layout :Grid{Columns:2,Spacing:10},
					Children: []Widget{
						VSplitter{
							Children:	[]Widget{
								LineEdit{
									MinSize:	Size{Width: 120,Height: 10},
									AssignTo:	&sw.editFile,
									ReadOnly: true,
								},
								PushButton{
									Text:    "Open file",
									MinSize: Size{Width: 120,Height: 50},
									OnClicked: func(){
										if err:= sw.pbClicked(); err!=nil {
											log.Print(err)
										}									
									},
									AssignTo: &sw.fileButton,
								},
								PushButton{
									Text:    "Start Attack",
									MinSize: Size{Width: 120,Height: 50},
									OnClicked: func(){
										go ExcBrute(df.cookieDVWA,df.urlDVWA,sw)
										return
									},
									AssignTo: &sw.attackButton,
								},
								TextEdit{
									MinSize:	Size{Width: 120,Height: 10},
									AssignTo:	&sw.outPut,
									ReadOnly: true,
								},
							},
						},
					},
				},
			},
		}

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
								AssignTo:	&df.confirmButton,
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
											df.cookieDVWA = cookie
											df.urlDVWA = df.dvwaEdit.Text()
											df.urlLabel.SetText("PHPSSID:")
											df.dvwaEdit.SetText(cookie)
											df.dvwaEdit.SetReadOnly(true)
											df.confirmRes = "ok"
											return
										}
									}
								},
							},
							VSplitter{
								Children:	[]Widget{
									Label{
										AssignTo:	&df.itemInfo,
										Text:	"DVWA ITEM:",
										Visible: true,
									},
									PushButton{
										Text:    "Brute Force",
										MinSize: Size{Width: 120,Height: 50},
										OnClicked: func() {
											if df.confirmRes != "ok"{
												var tmp walk.Form
												walk.MsgBox(tmp, "Error", "Get cookie first!", walk.MsgBoxIconInformation)
												return
											}
											subdef.Create()
											//_ = sw.MainWindow.Run()
											return
										},
									},
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
	return df, nil
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
