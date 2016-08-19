/**
 * Copyright 2016 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package virtual

import (
	"time"

	"github.ibm.com/riethm/gopherlayer.git/datatypes"
	"github.ibm.com/riethm/gopherlayer.git/helpers/product"
	"github.ibm.com/riethm/gopherlayer.git/services"
	"github.ibm.com/riethm/gopherlayer.git/session"
)

// Upgrade a virtual guest specified with an id, and a set of features to
// upgrade. When the upgrade takes place can also be specified (`when`), but
// this is optional. The time set will be 'now' if left as nil.
// The features to upgrade are specified as the options used in
// GetProductPrices().
func UpgradeVirtualGuest(
	sess session.Session,
	id int,
	options map[int]string,
	when time.Time,
) (datatypes.Container_Product_Order_Receipt, error) {

	prices, err := product.GetProductPrices(sess, "VIRTUAL_SERVER_INSTANCE", options)
	if err != nil {
		return nil, err
	}

	if when == nil {
		when = time.Now().UTC().Format(time.RFC3339)
	}

	order := datatypes.Container_Product_Order{
		VirtualGuests: []datatypes.Virtual_Guest{
			{Id: id},
		},
		Prices: prices,
		Properties: []datatypes.Container_Product_Order_Property{
			{
				Name:  "MAINTENANCE_WINDOW",
				Value: when,
			},
		},
	}

	orderService := services.GetProductOrderService(sess)
	return orderService.PlaceOrder(order, false)
}
