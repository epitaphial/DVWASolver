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

//for brute and injection
type BruteSubWindow struct{
		//BruteSubWindow
		editFile	 *walk.LineEdit
		fileButton		*walk.PushButton
		attackButton		*walk.PushButton
		filePath 		string
		outPut	 *walk.TextEdit
		*walk.MainWindow
		pushAble bool
		progressBar		*walk.ProgressBar
}

func (brt_sw *BruteSubWindow) pbClicked()error {

    dlg := new(walk.FileDialog)
    dlg.FilePath = brt_sw.filePath
    dlg.Title = "Select File"
	dlg.Filter = "Txt files (*.txt)|*.txt|All files (*.*)|*.*"
	brt_sw.editFile.SetText("")
	if ok,err:= dlg.ShowOpen(brt_sw); err!=nil {
        return err
    } else if!ok {
        return nil
    }
	brt_sw.editFile.SetText(dlg.FilePath)
	brt_sw.filePath = dlg.FilePath
	return nil
}


func createDVWA() (*DVWAfker, error) {
	//main window
	df := &DVWAfker{}
	//for brute and injection
	brt_sw := &BruteSubWindow{}
	inj_sw := &BruteSubWindow{}


	//brute force bind
	brute_subdef := MainWindow{
		AssignTo: &brt_sw.MainWindow,
		Title:    "BruteForce",
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
									AssignTo:	&brt_sw.editFile,
									ReadOnly: true,
								},
								PushButton{
									Text:    "Open file",
									MinSize: Size{Width: 120,Height: 50},
									OnClicked: func(){
										if err:= brt_sw.pbClicked(); err!=nil {
											log.Print(err)
										}									
									},
									AssignTo: &brt_sw.fileButton,
								},
								PushButton{
									Text:    "Start Attack",
									MinSize: Size{Width: 120,Height: 50},
									OnClicked: func(){
										if brt_sw.filePath == "" {
											var tmp walk.Form
											walk.MsgBox(tmp, "Warning", "Load File First!", walk.MsgBoxIconInformation)
										}else{
											if brt_sw.pushAble == true{
												go ExcBrute(df.cookieDVWA,df.urlDVWA,brt_sw)
											}else{
												var tmp walk.Form
												walk.MsgBox(tmp, "Warning", "Wait For The Current Progress!", walk.MsgBoxIconInformation)
											}
											}
											
										return
									},
									AssignTo: &brt_sw.attackButton,
								},
								TextEdit{
									MinSize:	Size{Width: 120,Height: 10},
									AssignTo:	&brt_sw.outPut,
									ReadOnly: true,
								},
								ProgressBar{
									AssignTo:	&brt_sw.progressBar,
								},
							},
						},
					},
				},
			},
		}

		//Command Injection bind
		inj_subdef := MainWindow{
			AssignTo: &inj_sw.MainWindow,
			Title:    "Command Injection",
			MinSize:  Size{Width: 200, Height: 250},
			Size:     Size{Width: 400, Height: 500},
			Layout:	VBox{},
			Children:	[]Widget{
					Composite{
						Layout :Grid{Columns:2,Spacing:10},
						Children: []Widget{
							VSplitter{
								Children:	[]Widget{
									PushButton{
										Text:    "Start Attack",
										MinSize: Size{Width: 120,Height: 50},
										OnClicked: func(){
											if inj_sw.pushAble == true{
												go ExcComInj(df.cookieDVWA,df.urlDVWA,inj_sw)
											}else{
												var tmp walk.Form
												walk.MsgBox(tmp, "Warning", "Wait For The Current Progress!", walk.MsgBoxIconInformation)
											}
											return
										},
										AssignTo: &inj_sw.attackButton,
									},
									TextEdit{
										MinSize:	Size{Width: 120,Height: 10},
										AssignTo:	&inj_sw.outPut,
										ReadOnly: true,
									},
									ProgressBar{
										AssignTo:	&inj_sw.progressBar,
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
										urlToParse := df.dvwaEdit.Text()
										if urlToParse[len(urlToParse)-1:] != "/"{
											urlToParse = urlToParse + "/"
										}
										
										cookie := InitDvwaUrl(urlToParse)
										if cookie == ""{
											var tmp walk.Form
											walk.MsgBox(tmp, "Error", "Can not get cookie!", walk.MsgBoxIconInformation)
											return
										}else{
											var tmp walk.Form
											walk.MsgBox(tmp, "Info", "Success init!", walk.MsgBoxIconInformation)
											df.cookieDVWA = cookie
											df.urlDVWA = urlToParse
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
											brt_sw.pushAble = true
											brute_subdef.Create()
											//_ = brt_sw.MainWindow.Run()
											return
										},
									},
									PushButton{
										Text:    "Command Injection",
										MinSize: Size{Width: 120,Height: 50},
										OnClicked: func() {
											if df.confirmRes != "ok"{
												var tmp walk.Form
												walk.MsgBox(tmp, "Error", "Get cookie first!", walk.MsgBoxIconInformation)
												return
											}
											inj_sw.pushAble = true
											inj_subdef.Create()
											//_ = brt_sw.MainWindow.Run()
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
