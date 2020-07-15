package controllers

import (
	"dyc/internal/helper"
	"dyc/internal/logger"
	"dyc/internal/module/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var Util *_util

type _util struct{}

func (*_util) PreserveHost(ctx *gin.Context) {
	helper.Success(ctx, gin.H{
		"header": ctx.Request.Header,
		"host":   ctx.Request.Host,
	})
}

func (*_util) Weather(ctx *gin.Context) {
	str := `{"data":{"alarm":{},"forecast_1h":{"0":{"degree":"21","update_time":"20200620080000","weather":"阴","weather_code":"02","weather_short":"阴","wind_direction":"东南风","wind_power":"3"},"1":{"degree":"22","update_time":"20200620090000","weather":"阴","weather_code":"02","weather_short":"阴","wind_direction":"东南风","wind_power":"3"},"10":{"degree":"23","update_time":"20200620180000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"11":{"degree":"23","update_time":"20200620190000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"12":{"degree":"23","update_time":"20200620200000","weather":"阴","weather_code":"02","weather_short":"阴","wind_direction":"东南风","wind_power":"3"},"13":{"degree":"22","update_time":"20200620210000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"14":{"degree":"22","update_time":"20200620220000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"15":{"degree":"21","update_time":"20200620230000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"3"},"16":{"degree":"21","update_time":"20200621000000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"3"},"17":{"degree":"21","update_time":"20200621010000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"3"},"18":{"degree":"21","update_time":"20200621020000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"3"},"19":{"degree":"21","update_time":"20200621030000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"3"},"2":{"degree":"22","update_time":"20200620100000","weather":"阴","weather_code":"02","weather_short":"阴","wind_direction":"东南风","wind_power":"3"},"20":{"degree":"21","update_time":"20200621040000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"3"},"21":{"degree":"21","update_time":"20200621050000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"3"},"22":{"degree":"21","update_time":"20200621060000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"3"},"23":{"degree":"21","update_time":"20200621070000","weather":"中雨","weather_code":"08","weather_short":"中雨","wind_direction":"东风","wind_power":"4"},"24":{"degree":"22","update_time":"20200621080000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"4"},"25":{"degree":"22","update_time":"20200621090000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"4"},"26":{"degree":"23","update_time":"20200621100000","weather":"中雨","weather_code":"08","weather_short":"中雨","wind_direction":"东风","wind_power":"4"},"27":{"degree":"23","update_time":"20200621110000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"5"},"28":{"degree":"23","update_time":"20200621120000","weather":"中雨","weather_code":"08","weather_short":"中雨","wind_direction":"东风","wind_power":"5"},"29":{"degree":"23","update_time":"20200621130000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"5"},"3":{"degree":"23","update_time":"20200620110000","weather":"阴","weather_code":"02","weather_short":"阴","wind_direction":"东南风","wind_power":"3"},"30":{"degree":"23","update_time":"20200621140000","weather":"阴","weather_code":"02","weather_short":"阴","wind_direction":"东风","wind_power":"5"},"31":{"degree":"23","update_time":"20200621150000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"4"},"32":{"degree":"23","update_time":"20200621160000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"4"},"33":{"degree":"23","update_time":"20200621170000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"4"},"34":{"degree":"22","update_time":"20200621180000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"3"},"35":{"degree":"22","update_time":"20200621190000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"4"},"36":{"degree":"22","update_time":"20200621200000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"4"},"37":{"degree":"22","update_time":"20200621210000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东风","wind_power":"4"},"38":{"degree":"22","update_time":"20200621220000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"39":{"degree":"22","update_time":"20200621230000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"4"},"4":{"degree":"23","update_time":"20200620120000","weather":"阴","weather_code":"02","weather_short":"阴","wind_direction":"东南风","wind_power":"3"},"40":{"degree":"22","update_time":"20200622000000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"41":{"degree":"23","update_time":"20200622010000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"4"},"42":{"degree":"23","update_time":"20200622020000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"4"},"43":{"degree":"23","update_time":"20200622030000","weather":"中雨","weather_code":"08","weather_short":"中雨","wind_direction":"东南风","wind_power":"4"},"44":{"degree":"23","update_time":"20200622040000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"4"},"45":{"degree":"23","update_time":"20200622050000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"4"},"46":{"degree":"23","update_time":"20200622060000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"4"},"47":{"degree":"23","update_time":"20200622070000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"4"},"5":{"degree":"23","update_time":"20200620130000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"6":{"degree":"24","update_time":"20200620140000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"7":{"degree":"23","update_time":"20200620150000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"8":{"degree":"23","update_time":"20200620160000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"},"9":{"degree":"23","update_time":"20200620170000","weather":"小雨","weather_code":"07","weather_short":"小雨","wind_direction":"东南风","wind_power":"3"}},"forecast_24h":{"0":{"day_weather":"阴","day_weather_code":"02","day_weather_short":"阴","day_wind_direction":"西风","day_wind_direction_code":"6","day_wind_power":"3","day_wind_power_code":"0","max_degree":"28","min_degree":"21","night_weather":"小雨","night_weather_code":"07","night_weather_short":"小雨","night_wind_direction":"东南风","night_wind_direction_code":"3","night_wind_power":"3","night_wind_power_code":"0","time":"2020-06-19"},"1":{"day_weather":"小雨","day_weather_code":"07","day_weather_short":"小雨","day_wind_direction":"东南风","day_wind_direction_code":"3","day_wind_power":"3","day_wind_power_code":"0","max_degree":"24","min_degree":"21","night_weather":"小雨","night_weather_code":"07","night_weather_short":"小雨","night_wind_direction":"东风","night_wind_direction_code":"2","night_wind_power":"4","night_wind_power_code":"1","time":"2020-06-20"},"2":{"day_weather":"中雨","day_weather_code":"08","day_weather_short":"中雨","day_wind_direction":"东风","day_wind_direction_code":"2","day_wind_power":"5","day_wind_power_code":"2","max_degree":"24","min_degree":"22","night_weather":"中雨","night_weather_code":"08","night_weather_short":"中雨","night_wind_direction":"东南风","night_wind_direction_code":"3","night_wind_power":"4","night_wind_power_code":"1","time":"2020-06-21"},"3":{"day_weather":"中雨","day_weather_code":"08","day_weather_short":"中雨","day_wind_direction":"东南风","day_wind_direction_code":"3","day_wind_power":"5","day_wind_power_code":"2","max_degree":"25","min_degree":"23","night_weather":"阴","night_weather_code":"02","night_weather_short":"阴","night_wind_direction":"东南风","night_wind_direction_code":"3","night_wind_power":"5","night_wind_power_code":"2","time":"2020-06-22"},"4":{"day_weather":"阴","day_weather_code":"02","day_weather_short":"阴","day_wind_direction":"东风","day_wind_direction_code":"2","day_wind_power":"3","day_wind_power_code":"0","max_degree":"27","min_degree":"23","night_weather":"阴","night_weather_code":"02","night_weather_short":"阴","night_wind_direction":"东南风","night_wind_direction_code":"3","night_wind_power":"4","night_wind_power_code":"1","time":"2020-06-23"},"5":{"day_weather":"小雨","day_weather_code":"07","day_weather_short":"小雨","day_wind_direction":"西风","day_wind_direction_code":"6","day_wind_power":"4","day_wind_power_code":"1","max_degree":"29","min_degree":"23","night_weather":"小雨","night_weather_code":"07","night_weather_short":"小雨","night_wind_direction":"西风","night_wind_direction_code":"6","night_wind_power":"4","night_wind_power_code":"1","time":"2020-06-24"},"6":{"day_weather":"小雨","day_weather_code":"07","day_weather_short":"小雨","day_wind_direction":"西风","day_wind_direction_code":"6","day_wind_power":"3","day_wind_power_code":"0","max_degree":"25","min_degree":"23","night_weather":"中雨","night_weather_code":"08","night_weather_short":"中雨","night_wind_direction":"东南风","night_wind_direction_code":"3","night_wind_power":"5","night_wind_power_code":"2","time":"2020-06-25"},"7":{"day_weather":"多云","day_weather_code":"01","day_weather_short":"多云","day_wind_direction":"东风","day_wind_direction_code":"2","day_wind_power":"4","day_wind_power_code":"1","max_degree":"31","min_degree":"24","night_weather":"多云","night_weather_code":"01","night_weather_short":"多云","night_wind_direction":"东南风","night_wind_direction_code":"3","night_wind_power":"3","night_wind_power_code":"0","time":"2020-06-26"}},"index":{"airconditioner":{"detail":"您将感到很舒适，一般不需要开启空调。","info":"较少开启","name":"空调开启"},"allergy":{"detail":"天气条件不易诱发过敏，有降水，特殊体质人群应预防感冒可能引发的过敏。","info":"不易发","name":"过敏"},"carwash":{"detail":"不宜洗车，未来24小时内有雨，如果在此期间洗车，雨水和路上的泥水可能会再次弄脏您的爱车。","info":"不宜","name":"洗车"},"chill":{"detail":"温度未达到风寒所需的低温，稍作防寒准备即可。","info":"无","name":"风寒"},"clothes":{"detail":"建议着长袖T恤、衬衫加单裤等服装。年老体弱者宜着针织长袖衬衫、马甲和长裤。","info":"舒适","name":"穿衣"},"cold":{"detail":"天气转凉，空气湿度较大，较易发生感冒，体质较弱的朋友请注意适当防护。","info":"较易发","name":"感冒"},"comfort":{"detail":"白天温度适宜，风力不大，相信您在这样的天气条件下，应会感到比较清爽和舒适。","info":"舒适","name":"舒适度"},"diffusion":{"detail":"气象条件有利于空气污染物稀释、扩散和清除。","info":"良","name":"空气污染扩散条件"},"dry":{"detail":"有降水，路面潮湿，车辆易打滑，请小心驾驶。","info":"潮湿","name":"路况"},"drying":{"detail":"有降水，不适宜晾晒。若需要晾晒，请在室内准备出充足的空间。","info":"不宜","name":"晾晒"},"fish":{"detail":"天气不好，有风，不适合垂钓。","info":"不宜","name":"钓鱼"},"heatstroke":{"detail":"天气舒适，对易中暑人群来说非常友善。","info":"无中暑风险","name":"中暑"},"makeup":{"detail":"风力不大，建议用中性保湿型霜类化妆品，无需选用防晒化妆品。","info":"保湿","name":"化妆"},"mood":{"detail":"有降水，雨水可能会使心绪无端地挂上轻愁，与其因下雨而无精打采，不如放松心情，好好欣赏一下雨景。你会发现雨中的世界是那般洁净温和、清新葱郁。","info":"较差","name":"心情"},"morning":{"detail":"有较强降水，风力稍大，请避免户外晨练，建议在室内做适当锻炼，保持身体健康。","info":"不宜","name":"晨练"},"sports":{"detail":"有降水，推荐您在室内进行健身休闲运动；若坚持户外运动，须注意携带雨具并注意避雨防滑。","info":"较不宜","name":"运动"},"sunglasses":{"detail":"白天有降水天气，视线较差，不需要佩戴太阳镜","info":"不需要","name":"太阳镜"},"sunscreen":{"detail":"属弱紫外辐射天气，长期在户外，建议涂擦SPF在8-12之间的防晒护肤品。","info":"弱","name":"防晒"},"time":"20200620","tourism":{"detail":"温度适宜，又有较弱降水和微风作伴，会给您的旅行带来意想不到的景象，适宜旅游，可不要错过机会呦！","info":"适宜","name":"旅游"},"traffic":{"detail":"有降水，路面湿滑，刹车距离延长，事故易发期，注意车距，务必小心驾驶。","info":"一般","name":"交通"},"ultraviolet":{"detail":"属弱紫外线辐射天气，无需特别防护。若长期在户外，建议涂擦SPF在8-12之间的防晒护肤品。","info":"最弱","name":"紫外线强度"},"umbrella":{"detail":"有降水，请带上雨伞，如果你喜欢雨中漫步，享受大自然给予的温馨和快乐，在短时间外出可收起雨伞。","info":"带伞","name":"雨伞"}},"limit":{"tail_number":"","time":""},"observe":{"degree":"23","humidity":"85","precipitation":"0.1","pressure":"1009","update_time":"202006200840","weather":"雨","weather_code":"301","weather_short":"雨","wind_direction":"8","wind_power":"0"},"rise":{"0":{"sunrise":"04:50","sunset":"19:01","time":"20200620"},"1":{"sunrise":"04:50","sunset":"19:01","time":"20200621"},"10":{"sunrise":"04:53","sunset":"19:02","time":"20200630"},"11":{"sunrise":"04:53","sunset":"19:02","time":"20200701"},"12":{"sunrise":"04:54","sunset":"19:02","time":"20200702"},"13":{"sunrise":"04:54","sunset":"19:02","time":"20200703"},"14":{"sunrise":"04:55","sunset":"19:02","time":"20200704"},"2":{"sunrise":"04:50","sunset":"19:01","time":"20200622"},"3":{"sunrise":"04:51","sunset":"19:01","time":"20200623"},"4":{"sunrise":"04:51","sunset":"19:01","time":"20200624"},"5":{"sunrise":"04:51","sunset":"19:02","time":"20200625"},"6":{"sunrise":"04:51","sunset":"19:02","time":"20200626"},"7":{"sunrise":"04:52","sunset":"19:02","time":"20200627"},"8":{"sunrise":"04:52","sunset":"19:02","time":"20200628"},"9":{"sunrise":"04:52","sunset":"19:02","time":"20200629"}},"tips":{"forecast_24h":{"0":"今天有小雨，出门记得带伞~"},"observe":{"0":"下雨了，出门记得带伞~","1":"现在的温度比较舒适~"}}},"message":"OK","status":200}`
	var data interface{}
	_ = json.Unmarshal([]byte(str), &data)
	ctx.JSONP(200, data)
}

func (*_util) Ip(ctx *gin.Context) {
	ip := ctx.ClientIP()
	//ip := "175.44.108.169"
	res, err := util.LocationByIp(ip)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, res)
	return
}

func (*_util) Amap(ctx *gin.Context) {
	latitude, exists := ctx.GetQuery("latitude")
	if !exists {
		helper.Fail(ctx, errors.New("请指定经纬度"))
		return
	}
	longitude, exists := ctx.GetQuery("longitude")
	if !exists {
		helper.Fail(ctx, errors.New("请指定经纬度"))
		return
	}
	if address, err := util.Location.FindByGeoCode(ctx, latitude, longitude); err != nil {
		helper.Fail(ctx, err)
		return
	} else {
		helper.Success(ctx, address)
		return
	}
}

func (*_util) Location(ctx *gin.Context) {
	latitude := ctx.Query("latitude")
	longitude := ctx.Query("longitude")
	if latitude != "" && longitude != "" {
		if address, err := util.Location.FindByGeoCode(ctx, latitude, longitude); err == nil {
			if address.AdCode != "" {
				if res, err := util.AdCoder.Component(ctx, string(address.AdCode)); err == nil {
					helper.Success(ctx, res)
					return
				}
			}
		} else {
			logger.Debugf("amap request error: %s", err)
		}
	}
	ip := ctx.ClientIP()
	if l, err := util.LocationByIp(ip); err != nil {
		helper.Fail(ctx, err)
		return
	} else {
		if l["city"] != "" {
			if city, err := util.AdCoder.FindCity(ctx, l["city"]); err != nil {
				helper.Fail(ctx, err)
				return
			} else {
				if res, err := util.AdCoder.Component(ctx, city.Adcode); err != nil {
					helper.Fail(ctx, err)
					return
				} else {
					helper.Success(ctx, res)
					return
				}
			}
		}
		defualt := gin.H{
			"city": gin.H{
				"name": "北京市",
				"adcode": "110100",
				"city_code": "",
			},
			"province": gin.H{
				"name": "北京市",
				"adcode": "110000",
				"city_code": "",
			},
		}
		helper.Success(ctx, defualt)
		return
	}
}
