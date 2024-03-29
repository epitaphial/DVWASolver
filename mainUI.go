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
		editCom		*walk.LineEdit
		comExc		*walk.PushButton
		fileButton		*walk.PushButton
		attackButton		*walk.PushButton
		filePath 		string
		outPut	 *walk.TextEdit
		*walk.MainWindow
		pushAble bool
		progressBar		*walk.ProgressBar
}

func (brt_sw *BruteSubWindow) pbClicked(pbTp int)error {

    dlg := new(walk.FileDialog)
    dlg.FilePath = brt_sw.filePath
	dlg.Title = "Select File"
	if pbTp == 0{
		dlg.Filter = "Txt files (*.txt)|*.txt|All files (*.*)|*.*"
	}else if pbTp == 1{
		dlg.Filter = "All files (*.*)|*.*"
	}

	
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
	file_sw := &BruteSubWindow{}
	xss_s_sw := &BruteSubWindow{}


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
										if err:= brt_sw.pbClicked(0); err!=nil {
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

		//file upload
		fileUp_subdef := MainWindow{
			AssignTo: &file_sw.MainWindow,
			Title:    "FileUpload",
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
										AssignTo:	&file_sw.editFile,
										ReadOnly: true,
									},
									PushButton{
										Text:    "Open file",
										MinSize: Size{Width: 120,Height: 50},
										OnClicked: func(){
											if err:= file_sw.pbClicked(1); err!=nil {
												log.Print(err)
											}									
										},
										AssignTo: &file_sw.fileButton,
									},
									PushButton{
										Text:    "Start Attack",
										MinSize: Size{Width: 120,Height: 50},
										OnClicked: func(){
											if file_sw.filePath == "" {
												var tmp walk.Form
												walk.MsgBox(tmp, "Warning", "Load File First!", walk.MsgBoxIconInformation)
											}else{
												if file_sw.pushAble == true{
													go ExcFileup(df.cookieDVWA,df.urlDVWA,file_sw)
												}else{
													var tmp walk.Form
													walk.MsgBox(tmp, "Warning", "Wait For The Current Progress!", walk.MsgBoxIconInformation)
												}
												}	
											return
										},
										AssignTo: &file_sw.attackButton,
									},
									LineEdit{
										AssignTo:	&file_sw.editCom,
									},
									PushButton{
										Text:    "Excute Command",
										MinSize: Size{Width: 120,Height: 50},
										OnClicked: func(){
											ExcFileCommand(df.urlDVWA,file_sw)
										},
										AssignTo: &file_sw.comExc,
									},
									TextEdit{
										MinSize:	Size{Width: 120,Height: 10},
										AssignTo:	&file_sw.outPut,
										ReadOnly: true,
									},
									ProgressBar{
										AssignTo:	&file_sw.progressBar,
									},
								},
							},
						},
					},
				},
			}

		xss_s_subdef := MainWindow{
			AssignTo: &xss_s_sw.MainWindow,
			Title:    "XSS Stored",
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
											if xss_s_sw.pushAble == true{
												go ExcXssStore(df.cookieDVWA,df.urlDVWA,xss_s_sw)
											}else{
												var tmp walk.Form
												walk.MsgBox(tmp, "Warning", "Wait For The Current Progress!", walk.MsgBoxIconInformation)
											}
											return
										},
										AssignTo: &xss_s_sw.attackButton,
									},
									TextEdit{
										MinSize:	Size{Width: 120,Height: 10},
										AssignTo:	&xss_s_sw.outPut,
										ReadOnly: true,
									},
									ProgressBar{
										AssignTo:	&xss_s_sw.progressBar,
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
		//Icon:     "./res/favicon.ico",
		MinSize:  Size{Width: 320, Height: 300},
		Size:     Size{Width: 640, Height: 600},
		Layout:   VBox{},
		MenuItems: []MenuItem{
			Menu{
				Text: "File",
				Items: []MenuItem{
					Action{
						Text: "Exit",
						OnTriggered: func() {
							df.window.Close()
						},
					},
				},
			},
			Menu{
				Text: "Help",
				Items: []MenuItem{
					Action{
						Text: "About",
						OnTriggered: func() {
							walk.MsgBox(df.window, "about", "all right reserved by Curled",
								walk.MsgBoxIconInformation|walk.MsgBoxDefButton1)
						},
					},
				},
			},
		},
		Children: []Widget{
			GroupBox{
                Layout: HBox{},
                Children: []Widget{
					Label{
						AssignTo:	&df.urlLabel,
						Text:	"DVWA URL:",
						Visible: true,
					},
                    LineEdit{
						MinSize:	Size{Width: 120,Height: 20},
                        AssignTo: &df.dvwaEdit,
                    },

					PushButton{
						Text:    "Get Cookie!",
						MinSize: Size{Width: 120,Height: 20},
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
                },
			},
			GroupBox{
                Layout: HBox{},
                Children: []Widget{
							VSplitter{
								Children:	[]Widget{
									Label{
										AssignTo:	&df.itemInfo,
										Text:	"DVWA ITEM:",
										Visible: true,
										TextAlignment: AlignCenter,
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
											return
										},
									},
									PushButton{
										Text:    "File Upload",
										MinSize: Size{Width: 120,Height: 50},
										OnClicked: func() {
											if df.confirmRes != "ok"{
												var tmp walk.Form
												walk.MsgBox(tmp, "Error", "Get cookie first!", walk.MsgBoxIconInformation)
												return
											}
											file_sw.pushAble = true
											fileUp_subdef.Create()
											return
										},
									},
									PushButton{
										Text:    "XSS(Stored)",
										MinSize: Size{Width: 120,Height: 50},
										OnClicked: func() {
											if df.confirmRes != "ok"{
												var tmp walk.Form
												walk.MsgBox(tmp, "Error", "Get cookie first!", walk.MsgBoxIconInformation)
												return
											}
											xss_s_sw.pushAble = true
											xss_s_subdef.Create()
											return
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
