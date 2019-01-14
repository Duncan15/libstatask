package work

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"golang.org/x/net/html"
	"libstatask/common/dbs"
	"libstatask/common/nets"
	"libstatask/common/sys"
	"net/http"
	"net/url"
	"time"
)

var (
	nameHandlerMap = map[string]func(node *html.Node) []*dbs.LibUnit{
		"空间": func(node *html.Node) []*dbs.LibUnit {
			return nil
		},
		"座位": func(node *html.Node) []*dbs.LibUnit {
			libUnitSlice := []*dbs.LibUnit{}
			se := goquery.NewDocumentFromNode(node)
			se.Find("li.it").Each(func(i int, selection *goquery.Selection) {
				unit := new(dbs.LibUnit)
				unit.Name = selection.Find("span").Get(0).FirstChild.Data
				for _, v := range selection.Get(0).Attr {
					if v.Key == "url" {
						uri, _ := url.Parse(v.Val)
						unit.RoomID = uri.Query().Get("roomId")
					}
				}
				libUnitSlice = append(libUnitSlice, unit)
			})
			return libUnitSlice
		},
	}
)

//CollectLibDomainLink collect domain info
func CollectLibDomainLink() {
	//store the address and param pair
	libUnitSlice := []*dbs.LibUnit{}
	if sys.IsRoot() && !nets.EasyPing("10.12.162.31") {
		glog.Warningf("the machine is not in the school nets, exit this task")
		return
	}
	homePage := "http://10.12.162.31/ClientWeb/xcus/ic2/Default.aspx"
	resp, err := http.Get(homePage)
	if err != nil {
		glog.Errorf("error happen when get homepage %s, the error content is %v", homePage, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		glog.Warningf("the homepage response code is %d", resp.StatusCode)
		return
	}

	//load the html document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		glog.Errorf("parse document error: %v", err)
		return
	}

	//the target list
	se := doc.Find("#info_tree #item_list .it_cls_list").Children()
	fun := func(node *html.Node) []*dbs.LibUnit { return nil }
	for i := 0; i < se.Size(); i++ {
		if f, has := nameHandlerMap[se.Get(i).FirstChild.Data]; has {
			fun = f
			continue
		}
		libUnitSlice = append(libUnitSlice, fun(se.Get(i))...)
	}

	session := dbs.MySQL.NewSession()
	defer session.Close()
	for _, unit := range libUnitSlice {
		unit.UpdateTime = time.Now().Unix()
		if has, err := session.Where("deleted=?", dbs.Undeleted).Exist(&dbs.LibUnit{Name: unit.Name}); err == nil && has {
			session.Where("name = ?", unit.Name).Update(unit)
		} else {
			session.InsertOne(unit)
		}
	}
}
