package main
/*
分布式事务
*/


import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

type (
	User struct {
		Id int64 `json:"id"`
		Username string `json:"username" xorm:"varchar(100) notnull 'username'`
		Password string `json:"password" xorm:"varchar(100) notnull 'password'`
	}
	Balance struct {
		Id int64 `json:"id"`
		UserId int64 `json:"userId" xorm:"int(11) notnull 'user_id'`
		Balance float64 `json:"balance" xorm:"double notnull 'balance'`
	}
	Goods struct {
		Id int64 `json:"id"`
		GoodsName string `json:"goodsName" xorm:"varchar(10) notnull 'goods_name'`
		Price float64 `json:"price" xorm:"double notnull 'price'`
		Stock int64 `json:"stock" xorm:"int(11) notnull 'stock'`
	}

	Order struct {
		Id int64 `json:"id"`
		UserId int64 `json:"userId" xorm:"int(11) notnull 'user_id'`
		Amount float64 `json:"amount" xorm:"double notnull 'amount'`
	}
)


func (b *Balance)UpdateBalanceByUserId(userId int64, amount float64,opt func()error) error{
	_, err := engineBalance.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		one, _ := b.findBalanceByUserId(userId)

		if i, err := session.Where("user_id = ?", userId).Update(&Balance{Balance: one.Balance - amount}); err != nil {
			log.Println(err)
			return i, err
		}

		return nil, opt()
	})
	return err
}


func (b *Balance)findBalanceByUserId(userId int64) (*Balance, error){
	balance := new(Balance)
	if _, err := engineBalance.Where("user_id = ?", userId).Get(balance); err != nil {
		log.Println(err)
		return nil, err
	}

	return balance, nil
}


func (g *Goods)UpdateStockByGoodsId(goodsId int64, opt func()error) error{
	_, err := engineGoods.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		one := g.findGoodsByGoodsId(goodsId)
		if i, err := session.Where("id = ?", goodsId).Update(&Goods{Stock:one.Stock - 2}); err != nil {
			log.Println(err)
			return i, err
		}

		return nil, opt()
	})
	return err
}


func (g *Goods)findGoodsByGoodsId(goodsId int64) *Goods {
	goods := new(Goods)
	if _, err := engineGoods.ID(goodsId).Get(goods); err != nil {
		log.Println(err)
		return nil
	}

	return goods
}


func (o *Order)InsertOrderRecord(userId int64, amount float64 ,opt func()error) error{
	_, err := engineOrder.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		if i, err := session.Insert(&Order{UserId: userId, Amount: amount}); err != nil {
			log.Println(err)
			return i, err
		}
		return nil, opt()
	})

	return err
}





var engineUser *xorm.Engine
var engineBalance *xorm.Engine
var engineGoods *xorm.Engine
var engineOrder *xorm.Engine
func main() {
	var err error
	engineUser, err = xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/members?charset=utf8")
	engineBalance, err = xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3307)/balance?charset=utf8")
	engineGoods, err = xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3308)/goods?charset=utf8")
	engineOrder, err = xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/members?charset=utf8")


	if err != nil {
		log.Fatal(err)
	}


	user := new(User)
	order := new(Order)
	balance := new(Balance)
	goods := new(Goods)

	if _, err := engineUser.Where("id = ?", 1).Get(user); err != nil {
		log.Fatal(err)
	}


	//订单服务
	orderNum := float64(2)
	goodsId := int64(1)
	goods = goods.findGoodsByGoodsId(goodsId)
	amount := orderNum * goods.Price

	err = order.InsertOrderRecord(user.Id, amount, func() error {
		//订单业务1
		//订单业务2
		//订单业务3
		//订单业务n
		return balance.UpdateBalanceByUserId(user.Id, amount,func() error {
			//余额服务
			//余额业务1
			//余额业务2
			//余额业务3
			//余额业务n
			return goods.UpdateStockByGoodsId(goods.Id, func() error {
				//商品业务1
				//商品业务2
				//商品业务3
				//商品业务n
				return nil
			})
		})
	})


	if err != nil {
		log.Fatal(err)
	}
	log.Println("success")




}
