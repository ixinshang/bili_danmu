package reply

import (
	"container/list"
	"errors"
	"time"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/gdk"
	p "github.com/qydysky/part"
	F "github.com/qydysky/bili_danmu/F"
)
const max = 50
const appId = "com.github.qydysky.bili_danmu.reply"

type gtk_list struct {
	text *gtk.TextView
	img *gtk.Image
	handle glib.SignalHandle
}
var pro_style *gtk.CssProvider
var gtkGetList = list.New()

var imgbuf = make(map[string](*gdk.Pixbuf))
var keep int
var keep_s int
var keep_key = map[string]int{
	"face/0default":0,
	"face/0room":0,
	"face/0buyguide":8,
	"face/0gift":8,
	"face/0jiezou":8,
	"face/0level1":5,
	"face/0level2":5,
	"face/0level3":5,
	"face/0superchat":13,
}
var (
	Gtk_on bool
	Gtk_Tra bool
	Gtk_img_path string = "face"
	Gtk_danmuChan chan string = make(chan string, 1000)
	Gtk_danmuChan_uid chan string = make(chan string, 1000)
)

func Gtk_danmu() {
	if Gtk_on {return}

	var y func(string,string)
	var win *gtk.Window
	var scrolledwindow0 *gtk.ScrolledWindow
	var viewport0 *gtk.Viewport

	application, err := gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	if err != nil {log.Println(err);return}

	application.Connect("startup", func() {
		log.Println("application startup")	
		var grid0 *gtk.Grid;

		builder, err := gtk.BuilderNewFromFile("ui/1.glade")
		if err != nil {log.Println(err);return}

		{
			signals := map[string]interface{}{
				"on_main_window_destroy": onMainWindowDestroy,
			}
			builder.ConnectSignals(signals)
		}
		{
			obj, err := builder.GetObject("main_window")
			if err != nil {log.Println(err);return}
			win, err = isWindow(obj)
			if err != nil {log.Println(err);return}
			application.AddWindow(win)
			defer win.ShowAll()
		}
		{
			obj, err := builder.GetObject("scrolledwindow0")
			if err != nil {log.Println(err);return}
			if tmp,ok := obj.(*gtk.ScrolledWindow); ok {
				scrolledwindow0 = tmp
			}else{log.Println("cant find #scrolledwindow0 in .glade");return}
		}

		{
			obj, err := builder.GetObject("viewport0")
			if err != nil {log.Println(err);return}
			if tmp,ok := obj.(*gtk.Viewport); ok {
				viewport0 = tmp
			}else{log.Println("cant find #viewport0 in .glade");return}
		}

		{
			obj, err := builder.GetObject("grid0")
			if err != nil {log.Println(err);return}
			if tmp,ok := obj.(*gtk.Grid); ok {
				grid0 = tmp
			}else{log.Println("cant find #grid0 in .glade");return}
		}

		imgbuf["face/0default"],_ = gdk.PixbufNewFromFileAtSize("face/0default", 40, 40);

		{
			var e error
			if pro_style,e = gtk.CssProviderNew();e == nil{
				if e = pro_style.LoadFromPath(`ui/1.css`);e != nil{
					log.Println(e)
				}
			}else{log.Println(e)}
		}

		y = func(s,img_src string){
			var tmp_list gtk_list

			tmp_list.text,_ = gtk.TextViewNew();
			{
				tmp_list.text.SetMarginStart(5)
				tmp_list.text.SetEditable(false)
				tmp_list.text.SetHExpand(true)
				tmp_list.text.SetWrapMode(gtk.WRAP_WORD_CHAR)
			}
			{
				var e error
				tmp_list.handle,e = tmp_list.text.Connect("size-allocate", func(){

					b,e := tmp_list.text.GetBuffer()
					if e != nil {log.Println(e);return}
					b.SetText(s)

					{
						var e error
						tmp := scrolledwindow0.GetVAdjustment()
						h := viewport0.GetViewWindow().WindowGetHeight()
						if tmp.GetUpper() - tmp.GetValue() < float64(h) * 1.7 {
							tmp.SetValue(tmp.GetUpper() - float64(h))
						}
						if e != nil {log.Println(e)}
					}
				})
				if e != nil {log.Println(e)}
			}

			tmp_list.img,_ =gtk.ImageNew();
			{
				var (
					pixbuf *gdk.Pixbuf
					e error
				)
				if v,ok := imgbuf[img_src];ok{
					pixbuf,e = gdk.PixbufCopy(v)
				} else {
					pixbuf,e = gdk.PixbufNewFromFileAtSize(img_src, 40, 40);
					if e == nil {imgbuf[img_src],e = gdk.PixbufCopy(pixbuf)}
				}
				if e == nil {tmp_list.img.SetFromPixbuf(pixbuf)}
			}
			{
				loc := int(grid0.Container.GetChildren().Length())/2;

				if sec,ok := keep_key[img_src];ok && sec != 0 {
					if sty,e := tmp_list.text.GetStyleContext();e == nil{
						sty.AddClass("view")
						sty.AddProvider(pro_style,gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
					}else{{log.Println(e)}}
					
					grid0.InsertRow(loc);
					grid0.Attach(tmp_list.img, 0, loc, 1, 1)
					grid0.Attach(tmp_list.text, 1, loc, 1, 1)
					keep += 1
					keep_s += sec
					if keep_s > 20 {keep_s = 20}
					
					glib.TimeoutAdd(uint(keep_s * 1000),func()(o bool){
						o = false
						keep -= 1
						return
					})
				}else{
					grid0.InsertRow(loc - keep);
					grid0.Attach(tmp_list.img, 0, loc - keep, 1, 1)
					grid0.Attach(tmp_list.text, 1, loc - keep, 1, 1)
				}

				loc = int(grid0.Container.GetChildren().Length())/2;
				for loc > max {
					if i,e := grid0.GetChildAt(0,0); e != nil{i.(*gtk.Widget).Destroy()}
					if i,e := grid0.GetChildAt(1,0); e != nil{i.(*gtk.Widget).Destroy()}
					grid0.RemoveRow(0)
					loc -= 1
				}
			}
			
			go func(){
				for{
					time.Sleep(time.Second)
					if keep_s > 0 {keep_s -= 1}
				}
			}()

			glib.TimeoutAdd(uint(3000), func()(o bool){
				o = true

				if gtkGetList.Len() == 0 {return}
				el := gtkGetList.Front()
				if el == nil {return}
				if uid,ok := gtkGetList.Remove(el).(string);ok{
					go func(){
						src := F.Get_face_src(uid)
						if src == "" {return}
						req := p.Req()
						if e := req.Reqf(p.Rval{
							Url:src,
							SaveToPath:Gtk_img_path + `/` + uid,
							Timeout:3,
						}); e != nil{log.Println(e);}
					}()
				}

				{
					if len(imgbuf) > 1000 {
						for k,_ := range imgbuf {delete(imgbuf,k);break}
					}
				}
				return
			})

			win.ShowAll()
		}

		Gtk_on = true
	})

	application.Connect("activate", func() {
		log.Println("application activate")
		glib.TimeoutAdd(uint(300),func()(o bool){
			o = true
			var tmax int = max
			for len(Gtk_danmuChan) != 0 {
				tmax -= 1
				if tmax <= 0 {return}
				y(<-Gtk_danmuChan,load_face(<-Gtk_danmuChan_uid))
			}
			return
		})

	})

	application.Connect("shutdown", func() {
		log.Println("application shutdown")	
		Gtk_on = false
	})

	application.Run(nil)
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func onMainWindowDestroy() {
	log.Println("onMainWindowDestroy")
}

func load_face(uid string) (loc string) {
	loc = Gtk_img_path + `/` + "0default"
	if uid == "" {return}
	if _,ok := keep_key[uid];ok{
		loc = Gtk_img_path + `/` + uid
		return
	}
	if p.Checkfile().IsExist(Gtk_img_path + `/` + uid) && p.Rand().MixRandom(1,100) > 1 {
		loc = Gtk_img_path + `/` + uid
		return
	}
	if gtkGetList.Len() > 1000 {return}
	gtkGetList.PushBack(uid)
	return
}