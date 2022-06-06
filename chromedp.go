package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"time"
)

func main() {
	fmt.Println("Login")
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	//Setting timeout
	//ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()

	var nodesButtonList []*cdp.Node
	//selector for checking the Data Services button
	selector := "#page-sidebar div ul li button"

	url := "https://console-openshift-console.apps.rhoda-lab.51ty.p1.openshiftapps.com/dashboards"

	if err := chromedp.Run(ctx,
		SetCookie("openshift-session-token", "", "console-openshift-console.apps.rhoda-lab.51ty.p1.openshiftapps.com", "/", false, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#page-sidebar`),
		chromedp.Nodes(selector, &nodesButtonList),
	); err != nil {
		panic(err)
	}

	selectorLi := "section ul li a"
	var dataServiceNode []*cdp.Node
	//Get Data Services button
	buttonLi := getLi(nodesButtonList)
	if err := chromedp.Run(ctx,
		chromedp.Nodes(selectorLi, &dataServiceNode, chromedp.ByQueryAll, chromedp.FromNode(buttonLi)),
	); err != nil {
		panic(err)
	}

	href := getHref(dataServiceNode)
	fmt.Println(href)
	var textExists bool
	u := fmt.Sprintf("https://console-openshift-console.apps.rhoda-lab.51ty.p1.openshiftapps.com%s", href)
	//	u := "https://console-openshift-console.apps.rhoda-lab.51ty.p1.openshiftapps.com" + href
	fmt.Printf(u)
	selectorDA := "#content-scrollable h1 div span"
	var dataAccessNode []*cdp.Node
	var outer string
	err := chromedp.Run(ctx,
		chromedp.Navigate(u),
		chromedp.WaitVisible(`#content-scrollable button`),
		//	chromedp.InnerHTML(selectorDA, &outer, chromedp.ByQuery),
		chromedp.Nodes(selectorDA, &dataAccessNode, chromedp.ByQuery),
		//chromedp.EvaluateAsDevTools(`document.querySelector("#content-scrollable h1 div" ).innerHTML.includes("Database Access")`, &textExists),
	)
	if err != nil {
		fmt.Println("Error", err.Error())
		panic(err.Error())
	}

	fmt.Println(outer)
	for _, nodeda := range dataAccessNode {
		//NodeName is the button here, looping through buttons to get Data services
		foundText := nodeda.Parent.Children[0].NodeValue
		fmt.Println(foundText)
	}
	//	status2 := resp.Status
	fmt.Println(textExists)
	if textExists == false {
		fmt.Println("Error")
	} else {
		fmt.Println("Found Database Access")
	}
}
func getHref(dataServiceNode []*cdp.Node) string {
	for _, aNode := range dataServiceNode {
		text := aNode.Children[0].NodeValue
		fmt.Println(text)
		if aNode.Children[0].NodeValue == "Database Access" {
			href := aNode.AttributeValue("href")
			return href
		}
	}
	return ""
}
func checkAdminDashboard(li *cdp.Node, ctx context.Context) {
	fmt.Println("checkAdminDashboard")
	fmt.Println(li)
	//	var data string
	selector := "section ul li a"
	var result []*cdp.Node

	if err := chromedp.Run(ctx,
		chromedp.Nodes(selector, &result, chromedp.ByQueryAll, chromedp.FromNode(li)),
	); err != nil {
		panic(err)
	}
	fmt.Println(result)
	for _, aNode := range result {
		//temp stuff for visibility
		text := aNode.Children[0].NodeValue
		fmt.Println(text)
		u := aNode.AttributeValue("href")
		fmt.Printf("node: %s | href = %s\n", aNode.LocalName, u)
		textSelector := "#content-scrollable h1 "
		var dataAccessResult []*cdp.Node

		if aNode.Children[0].NodeValue == "Database Access" {
			u := "https://console-openshift-console.apps.rhoda-lab.51ty.p1.openshiftapps.com" + aNode.AttributeValue("href")
			fmt.Printf("node: %s | href = %s\n", aNode.LocalName, u)
			var textExists bool

			err := chromedp.Run(ctx,
				chromedp.Navigate(u),
				chromedp.WaitVisible(`#content-scrollable button`),
				//	chromedp.OuterHTML("html", &data, chromedp.ByQuery),
				chromedp.Nodes(textSelector, &dataAccessResult))
			//chromedp.EvaluateAsDevTools(`document.querySelector("#content-scrollable h1 div" ).innerHTML.includes("Database Access")`, &textExists))
			if err != nil {
				fmt.Println("Error", err.Error())
				panic(err.Error())
			}
			//	status2 := resp.Status
			fmt.Println("second status code:", textSelector)
			fmt.Println(textExists)
			if textExists == false {
				fmt.Println("Error")
			} else {
				fmt.Println("Found Database Access")
			}
			//}
		}
	}
}

func SetCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		success := network.SetCookie(name, value).
			WithExpires(&expr).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSecure(secure).
			Do(ctx)
		if success != nil {
			return fmt.Errorf("could not set cookie %s", name)
		}
		return nil
	})
}

func getLi(nodesButtonList []*cdp.Node) *cdp.Node {
	for _, node := range nodesButtonList {
		//NodeName is the button here, looping through buttons to get Data services
		fmt.Println(node.NodeName)

		for _, child := range node.Children {
			//getting Data Services Button
			fmt.Println(child.NodeValue)
			if child.NodeValue == "Data Services" {
				//get the parent's parent which is Li to click on the Database Access button
				return child.Parent.Parent
			}
		}
	}
	return nil
}
