package main

/*
 Six910 is a shopping cart and E-commerce system.
 Copyright (C) 2020 Ulbora Labs LLC. (www.ulboralabs.com)
 All rights reserved.

 Copyright (C) 2020 Ken Williamson
 All rights reserved.

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

*/

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	px "github.com/Ulbora/GoProxy"
	lg "github.com/Ulbora/Level_Logger"

	hand "github.com/Ulbora/Six910/handlers"
	man "github.com/Ulbora/Six910/managers"
	db "github.com/Ulbora/dbinterface"
	mdb "github.com/Ulbora/dbinterface_mysql"
	sixmdb "github.com/Ulbora/six910-mysql"
	"github.com/gorilla/mux"

	jv "github.com/Ulbora/GoAuth2JwtValidator"

	_ "github.com/Ulbora/Six910/docs" // docs is generated by Swag CLI, you have to import it.

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Six910 API
// @version 1.0
// @description This is the Six910 (six nine ten) API
// @termsOfService https://github.com/Ulbora/Six910/blob/master/tos.html
// @contact.name API Support
// @contact.url http://www.ulboralabs.com/contact/form
// @license.name GPL-3.0
// @license.url https://github.com/Ulbora/Six910/blob/master/LICENSE
// @host localhost:3002
// @BasePath /
func main() {

	var mydb mdb.MyDB
	var proxy px.GoProxy
	var dbHost string
	var dbUser string
	var dbPassword string
	var dbName string

	var apiKey string

	var userHost string

	var l lg.Logger
	l.LogLevel = lg.AllLevel

	if os.Getenv("SIX910_DB_HOST") != "" {
		dbHost = os.Getenv("SIX910_DB_HOST")
	} else {
		dbHost = "localhost:3306"
	}

	if os.Getenv("SIX910_DB_USER") != "" {
		dbUser = os.Getenv("SIX910_DB_USER")
	} else {
		dbUser = "admin"
	}

	if os.Getenv("SIX910_DB_PASSWORD") != "" {
		dbPassword = os.Getenv("SIX910_DB_PASSWORD")
	} else {
		dbPassword = "admin"
	}

	if os.Getenv("SIX910_DB_DATABASE") != "" {
		dbName = os.Getenv("SIX910_DB_DATABASE")
	} else {
		dbName = "six910"
	}

	if os.Getenv("SIX910_USER_HOST") != "" {
		userHost = os.Getenv("SIX910_USER_HOST")
	} else {
		userHost = "http://localhost:3001"
	}

	if os.Getenv("SIX910_API_KEY") != "" {
		apiKey = os.Getenv("SIX910_API_KEY")
	} else {
		apiKey = "GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5"
	}

	mydb.Host = dbHost         // "localhost:3306"
	mydb.User = dbUser         // "admin"
	mydb.Password = dbPassword // "admin"
	mydb.Database = dbName     // "six910"
	var dbi db.Database = &mydb

	var sdb sixmdb.Six910Mysql
	//var l lg.Logger
	l.LogLevel = lg.AllLevel
	sdb.Log = &l
	sdb.DB = dbi
	dbi.Connect()

	var sm man.Six910Manager
	sm.Db = sdb.GetNew()
	sm.Log = &l
	sm.Proxy = proxy.GetNewProxy()
	sm.UserHost = userHost

	var sh hand.Six910Handler
	sh.Manager = sm.GetNew()
	sh.APIKey = apiKey
	sh.Log = &l

	var mc jv.MockOauthClient
	//mc.MockValidate = true
	sh.ValidatorClient = mc.GetNewClient()

	router := mux.NewRouter()
	port := "3002"
	envPort := os.Getenv("PORT")
	if envPort != "" {
		portInt, _ := strconv.Atoi(envPort)
		if portInt != 0 {
			port = envPort
		}
	}

	h := sh.GetNew()

	//sdb.MockAddStoreSuccess = true
	//sdb.MockStoreID = 5

	//h := sh.GetNew()

	var locacc man.LocalStoreAdminUser
	locacc.Username = "admin"
	locacc.Password = "admin"

	lstoreRes := sm.CreateLocalStore(&locacc)
	sm.Log.Debug("Creating local store", *lstoreRes)

	//store
	router.HandleFunc("/rs/store/add", h.AddStore).Methods("POST")
	router.HandleFunc("/rs/store/update", h.UpdateStore).Methods("PUT")
	router.HandleFunc("/rs/store/get/{storeName}/{localDomain}", h.GetStore).Methods("GET")
	router.HandleFunc("/rs/store/delete/{storeName}/{localDomain}", h.DeleteStore).Methods("DELETE")

	//customer
	router.HandleFunc("/rs/customer/add", h.AddCustomer).Methods("POST")
	router.HandleFunc("/rs/customer/update", h.UpdateCustomer).Methods("PUT")
	router.HandleFunc("/rs/customer/get/email/{email}/{storeId}", h.GetCustomer).Methods("GET")
	router.HandleFunc("/rs/customer/get/id/{id}/{storeId}", h.GetCustomerID).Methods("GET")
	router.HandleFunc("/rs/customer/get/list/{storeId}", h.GetCustomerList).Methods("GET")
	router.HandleFunc("/rs/customer/delete/{id}/{storeId}", h.DeleteCustomer).Methods("DELETE")

	//users
	router.HandleFunc("/rs/user/add", h.AddUser).Methods("POST")
	router.HandleFunc("/rs/user/update", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/rs/user/{username}/{storeId}", h.GetUser).Methods("GET")

	router.HandleFunc("/rs/user/get/admin/list/{storeId}", h.GetAdminUserList).Methods("GET")
	router.HandleFunc("/rs/user/get/customer/list/{storeId}", h.GetCustomerUserList).Methods("GET")

	//distributors
	router.HandleFunc("/rs/distributor/add", h.AddDistributor).Methods("POST")
	router.HandleFunc("/rs/distributor/update", h.UpdateDistributor).Methods("PUT")
	router.HandleFunc("/rs/distributor/get/id/{id}/{storeId}", h.GetDistributor).Methods("GET")
	router.HandleFunc("/rs/distributor/get/list/{storeId}", h.GetDistributorList).Methods("GET")
	router.HandleFunc("/rs/distributor/delete/{id}/{storeId}", h.DeleteDistributor).Methods("DELETE")

	//cart
	router.HandleFunc("/rs/cart/add", h.AddCart).Methods("POST")
	router.HandleFunc("/rs/cart/update", h.UpdateCart).Methods("PUT")
	router.HandleFunc("/rs/cart/get/{cid}/{storeId}", h.GetCart).Methods("GET")
	router.HandleFunc("/rs/cart/delete/{id}/{cid}/{storeId}", h.DeleteCart).Methods("DELETE")

	//cartItem
	router.HandleFunc("/rs/cartItem/add", h.AddCartItem).Methods("POST")
	router.HandleFunc("/rs/cartItem/update", h.UpdateCartItem).Methods("PUT")
	router.HandleFunc("/rs/cartItem/get/{cid}/{prodId}/{storeId}", h.GetCartItem).Methods("GET")
	router.HandleFunc("/rs/cartItem/get/list/{cartId}/{cid}/{storeId}", h.GetCartItemList).Methods("GET")
	router.HandleFunc("/rs/cartItem/delete/{id}/{prodId}/{cartId}", h.DeleteCartItem).Methods("DELETE")

	//address
	router.HandleFunc("/rs/address/add", h.AddAddress).Methods("POST")
	router.HandleFunc("/rs/address/update", h.UpdateAddress).Methods("PUT")
	router.HandleFunc("/rs/address/get/id/{id}/{cid}/{storeId}", h.GetAddress).Methods("GET")
	router.HandleFunc("/rs/address/get/list/{cid}/{storeId}", h.GetAddressList).Methods("GET")
	router.HandleFunc("/rs/address/delete/{id}/{cid}/{storeId}", h.DeleteAddress).Methods("DELETE")

	//category
	router.HandleFunc("/rs/category/add", h.AddCategory).Methods("POST")
	router.HandleFunc("/rs/category/update", h.UpdateCategory).Methods("PUT")
	router.HandleFunc("/rs/category/get/id/{id}/{storeId}", h.GetCategory).Methods("GET")
	router.HandleFunc("/rs/category/get/list/{storeId}", h.GetCategoryList).Methods("GET")
	router.HandleFunc("/rs/category/get/sub/list/{catId}", h.GetSubCategoryList).Methods("GET")
	router.HandleFunc("/rs/category/delete/{id}/{storeId}", h.DeleteCategory).Methods("DELETE")

	//shipping method
	router.HandleFunc("/rs/shippingMethod/add", h.AddShippingMethod).Methods("POST")
	router.HandleFunc("/rs/shippingMethod/update", h.UpdateShippingMethod).Methods("PUT")
	router.HandleFunc("/rs/shippingMethod/get/id/{id}/{storeId}", h.GetShippingMethod).Methods("GET")
	router.HandleFunc("/rs/shippingMethod/get/list/{storeId}", h.GetShippingMethodList).Methods("GET")
	router.HandleFunc("/rs/shippingMethod/delete/{id}/{storeId}", h.DeleteShippingMethod).Methods("DELETE")

	//shipping insurance
	router.HandleFunc("/rs/insurance/add", h.AddInsurance).Methods("POST")
	router.HandleFunc("/rs/insurance/update", h.UpdateInsurance).Methods("PUT")
	router.HandleFunc("/rs/insurance/get/id/{id}/{storeId}", h.GetInsurance).Methods("GET")
	router.HandleFunc("/rs/insurance/get/list/{storeId}", h.GetInsuranceList).Methods("GET")
	router.HandleFunc("/rs/insurance/delete/{id}/{storeId}", h.DeleteInsurance).Methods("DELETE")

	//product
	router.HandleFunc("/rs/product/add", h.AddProduct).Methods("POST")
	router.HandleFunc("/rs/product/update", h.UpdateProduct).Methods("PUT")
	router.HandleFunc("/rs/product/get/id/{id}/{storeId}", h.GetProductByID).Methods("GET")
	router.HandleFunc("/rs/product/get/sku/{sku}/{did}/{storeId}", h.GetProductBySku).Methods("GET")
	router.HandleFunc("/rs/product/get/promoted/{storeId}/{start}/{end}", h.GetProductsByPromoted).Methods("GET")
	router.HandleFunc("/rs/product/get/name/{name}/{storeId}/{start}/{end}", h.GetProductsByName).Methods("GET")
	router.HandleFunc("/rs/product/get/category/{catId}/{storeId}/{start}/{end}", h.GetProductsByCaterory).Methods("GET")
	router.HandleFunc("/rs/product/get/list/{storeId}/{start}/{end}", h.GetProductList).Methods("GET")
	router.HandleFunc("/rs/product/delete/{id}/{storeId}", h.DeleteProduct).Methods("DELETE")

	//Geographic Regions
	router.HandleFunc("/rs/region/add", h.AddRegion).Methods("POST")
	router.HandleFunc("/rs/region/update", h.UpdateRegion).Methods("PUT")
	router.HandleFunc("/rs/region/get/id/{id}/{storeId}", h.GetRegion).Methods("GET")
	router.HandleFunc("/rs/region/get/list/{storeId}", h.GetRegionList).Methods("GET")
	router.HandleFunc("/rs/region/delete/{id}/{storeId}", h.DeleteRegion).Methods("DELETE")

	//Geographic Sub Regions
	router.HandleFunc("/rs/subRegion/add", h.AddSubRegion).Methods("POST")
	router.HandleFunc("/rs/subRegion/update", h.UpdateSubRegion).Methods("PUT")
	router.HandleFunc("/rs/subRegion/get/id/{id}/{storeId}", h.GetSubRegion).Methods("GET")
	router.HandleFunc("/rs/subRegion/get/list/{regionId}/{storeId}", h.GetSubRegionList).Methods("GET")
	router.HandleFunc("/rs/subRegion/delete/{id}/{storeId}", h.DeleteSubRegion).Methods("DELETE")

	//excluded Geographic Sub Regions
	router.HandleFunc("/rs/excludedSubRegion/add", h.AddExcludedSubRegion).Methods("POST")
	router.HandleFunc("/rs/excludedSubRegion/get/list/{regionId}/{storeId}", h.GetExcludedSubRegionList).Methods("GET")
	router.HandleFunc("/rs/excludedSubRegion/delete/{id}/{regionId}/{storeId}", h.DeleteExcludedSubRegion).Methods("DELETE")

	//included Geographic Sub Regions
	router.HandleFunc("/rs/includedSubRegion/add", h.AddIncludedSubRegion).Methods("POST")
	router.HandleFunc("/rs/includedSubRegion/get/list/{regionId}/{storeId}", h.GetIncludedSubRegionList).Methods("GET")
	router.HandleFunc("/rs/includedSubRegion/delete/{id}/{regionId}/{storeId}", h.DeleteIncludedSubRegion).Methods("DELETE")

	//limit exclusions and inclusions to a zip code
	router.HandleFunc("/rs/zoneZip/add", h.AddZoneZip).Methods("POST")
	router.HandleFunc("/rs/zoneZip/exc/get/list/{exId}/{storeId}", h.GetZoneZipListByExclusion).Methods("GET")
	router.HandleFunc("/rs/zoneZip/inc/get/list/{incId}/{storeId}", h.GetZoneZipListByInclusion).Methods("GET")
	router.HandleFunc("/rs/zoneZip/delete/{id}/{incId}/{exId}/{storeId}", h.DeleteZoneZip).Methods("DELETE")

	//productCategory
	router.HandleFunc("/rs/productCategory/add", h.AddProductCategory).Methods("POST")
	router.HandleFunc("/rs/productCategory/delete/{categoryId}/{productId}/{storeId}", h.DeleteProductCategory).Methods("DELETE")

	//Orders
	router.HandleFunc("/rs/order/add", h.AddOrder).Methods("POST")
	router.HandleFunc("/rs/order/update", h.UpdateOrder).Methods("PUT")
	router.HandleFunc("/rs/order/get/id/{id}/{storeId}", h.GetOrder).Methods("GET")
	router.HandleFunc("/rs/order/get/list/{cid}/{storeId}", h.GetOrderList).Methods("GET")
	router.HandleFunc("/rs/order/get/store/list/{storeId}", h.GetStoreOrderList).Methods("GET")
	router.HandleFunc("/rs/order/get/store/list/status/{status}/{storeId}", h.GetStoreOrderListByStatus).Methods("GET")
	router.HandleFunc("/rs/order/delete/{id}/{storeId}", h.DeleteOrder).Methods("DELETE")

	//Order Items
	router.HandleFunc("/rs/orderItem/add", h.AddOrderItem).Methods("POST")
	router.HandleFunc("/rs/orderItem/update", h.UpdateOrderItem).Methods("PUT")
	router.HandleFunc("/rs/orderItem/get/id/{id}/{storeId}", h.GetOrderItem).Methods("GET")
	router.HandleFunc("/rs/orderItem/get/list/{orderId}/{storeId}", h.GetOrderItemList).Methods("GET")
	router.HandleFunc("/rs/orderItem/delete/{id}/{storeId}", h.DeleteOrderItem).Methods("DELETE")

	//Order Comments
	router.HandleFunc("/rs/orderComment/add", h.AddOrderComments).Methods("POST")
	router.HandleFunc("/rs/orderComment/get/list/{orderId}/{storeId}", h.GetOrderCommentList).Methods("GET")

	//Order Payment Transactions
	router.HandleFunc("/rs/orderTransaction/add", h.AddOrderTransaction).Methods("POST")
	router.HandleFunc("/rs/orderTransaction/get/list/{orderId}/{storeId}", h.GetOrderTransactionList).Methods("GET")

	//shipment
	router.HandleFunc("/rs/shipment/add", h.AddShipment).Methods("POST")
	router.HandleFunc("/rs/shipment/update", h.UpdateShipment).Methods("PUT")
	router.HandleFunc("/rs/shipment/get/id/{id}/{storeId}", h.GetShipment).Methods("GET")
	router.HandleFunc("/rs/shipment/get/list/{orderId}/{storeId}", h.GetShipmentList).Methods("GET")
	router.HandleFunc("/rs/shipment/delete/{id}/{storeId}", h.DeleteShipment).Methods("DELETE")

	//shipment boxes
	router.HandleFunc("/rs/shipmentBox/add", h.AddShipmentBox).Methods("POST")
	router.HandleFunc("/rs/shipmentBox/update", h.UpdateShipmentBox).Methods("PUT")
	router.HandleFunc("/rs/shipmentBox/get/id/{id}/{storeId}", h.GetShipmentBox).Methods("GET")
	router.HandleFunc("/rs/shipmentBox/get/list/{shipmentId}/{storeId}", h.GetShipmentBoxList).Methods("GET")
	router.HandleFunc("/rs/shipmentBox/delete/{id}/{storeId}", h.DeleteShipmentBox).Methods("DELETE")

	//Shipment Items in box
	router.HandleFunc("/rs/shipmentItem/add", h.AddShipmentItem).Methods("POST")
	router.HandleFunc("/rs/shipmentItem/update", h.UpdateShipmentItem).Methods("PUT")
	router.HandleFunc("/rs/shipmentItem/get/id/{id}/{storeId}", h.GetShipmentItem).Methods("GET")
	router.HandleFunc("/rs/shipmentItem/get/list/{shipmentId}/{storeId}", h.GetShipmentItemList).Methods("GET")
	router.HandleFunc("/rs/shipmentItem/get/list/box/{boxNumber}/{shipmentId}/{storeId}", h.GetShipmentItemListByBox).Methods("GET")
	router.HandleFunc("/rs/shipmentItem/delete/{id}/{storeId}", h.DeleteShipmentItem).Methods("DELETE")

	//Global Plugins
	router.HandleFunc("/rs/plugin/add", h.AddPlugin).Methods("POST")
	router.HandleFunc("/rs/plugin/update", h.UpdatePlugin).Methods("PUT")
	router.HandleFunc("/rs/plugin/get/id/{id}", h.GetPlugin).Methods("GET")
	router.HandleFunc("/rs/plugin/get/list/{start}/{end}", h.GetPluginList).Methods("GET")
	router.HandleFunc("/rs/plugin/delete/{id}", h.DeletePlugin).Methods("DELETE")

	//store plugins installed
	router.HandleFunc("/rs/storePlugin/add", h.AddStorePlugin).Methods("POST")
	router.HandleFunc("/rs/storePlugin/update", h.UpdateStorePlugin).Methods("PUT")
	router.HandleFunc("/rs/storePlugin/get/id/{id}/{storeId}", h.GetStorePlugin).Methods("GET")
	router.HandleFunc("/rs/storePlugin/get/list/{storeId}", h.GetStorePluginList).Methods("GET")
	router.HandleFunc("/rs/storePlugin/delete/{id}/{storeId}", h.DeleteStorePlugin).Methods("DELETE")

	//Plugins that are payment gateways
	router.HandleFunc("/rs/paymentGateway/add", h.AddPaymentGateway).Methods("POST")
	router.HandleFunc("/rs/paymentGateway/update", h.UpdatePaymentGateway).Methods("PUT")
	router.HandleFunc("/rs/paymentGateway/get/id/{id}/{storeId}", h.GetPaymentGateway).Methods("GET")
	router.HandleFunc("/rs/paymentGateway/get/list/{storeId}", h.GetPaymentGateways).Methods("GET")
	router.HandleFunc("/rs/paymentGateway/delete/{id}/{storeId}", h.DeletePaymentGateway).Methods("DELETE")

	//store shipment carrier like UPS and FEDex
	router.HandleFunc("/rs/shippingCarrier/add", h.AddShippingCarrier).Methods("POST")
	router.HandleFunc("/rs/shippingCarrier/update", h.UpdateShippingCarrier).Methods("PUT")
	router.HandleFunc("/rs/shippingCarrier/get/id/{id}/{storeId}", h.GetShippingCarrier).Methods("GET")
	router.HandleFunc("/rs/shippingCarrier/get/list/{storeId}", h.GetShippingCarrierList).Methods("GET")
	router.HandleFunc("/rs/shippingCarrier/delete/{id}/{storeId}", h.DeleteShippingCarrier).Methods("DELETE")

	//datastore------------------------------------
	router.HandleFunc("/rs/datastore/add", h.AddLocalDatastore).Methods("POST")
	router.HandleFunc("/rs/datastore/update", h.UpdateLocalDatastore).Methods("PUT")
	router.HandleFunc("/rs/datastore/get/{name}/{storeId}", h.GetLocalDatastore).Methods("GET")

	//instance--------------------
	router.HandleFunc("/rs/instance/add", h.AddInstance).Methods("POST")
	router.HandleFunc("/rs/instance/update", h.UpdateInstance).Methods("PUT")
	router.HandleFunc("/rs/instance/get/name/{name}/{dataStoreName}/{storeId}", h.GetInstance).Methods("GET")
	router.HandleFunc("/rs/instance/get/list/{dataStoreName}/{storeId}", h.GetInstanceList).Methods("GET")

	//write lock-------------
	router.HandleFunc("/rs/dataStoreWriteLock/add", h.AddDataStoreWriteLock).Methods("POST")
	router.HandleFunc("/rs/dataStoreWriteLock/update", h.UpdateDataStoreWriteLock).Methods("PUT")
	router.HandleFunc("/rs/dataStoreWriteLock/get/{dataStore}/{storeId}", h.GetDataStoreWriteLock).Methods("GET")

	fmt.Println("Six910 (six nine ten) server is running on port " + port + "!")

	// Swagger
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	http.ListenAndServe(":"+port, router)
}

// go mod init github.com/Ulbora/Six910
