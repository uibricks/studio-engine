// let exp = require('./expressions_utility')
// window.Eval = exp.evalExpression

const exp = require('./expressions_utility.go');

const md={
    "tYu6c8jd1BHprNcB4mztc": {
        "id": "tYu6c8jd1BHprNcB4mztc",
        "name": "GetAllNames",
        "raw": "[phd8JZrz5Z95uT95-Lmx7]",
        "refs": {
            "phd8JZrz5Z95uT95-Lmx7": {
                "path": [
                    {
                        "name": "data",
                        "type": "array"
                    },
                    {
                        "name": "name"
                    }
                ],
                "type": "stringArray"
            }
        },
        "type": "stringArray"
    }
}

const menu=[
    {
        "id": "tYu6c8jd1BHprNcB4mztc",
        "type": "stringArray"
    }
]

const data=`{"code":200,"meta":{"pagination":{"total":19,"pages":1,"page":1,"limit":20}},"data":[{"id":8,"name":"Gitanjali Verma","email":"gitanjali_verma@hayes.com","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:04.174+05:30","updated_at":"2021-05-31T03:50:04.174+05:30"},{"id":70,"name":"Deependra Verma","email":"verma_deependra@marks-fahey.name","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:05.220+05:30","updated_at":"2021-05-31T03:50:05.220+05:30"},{"id":97,"name":"Tapan Verma","email":"tapan_verma@bernhard.net","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:05.618+05:30","updated_at":"2021-05-31T03:50:05.618+05:30"},{"id":178,"name":"Adhrit Verma Sr.","email":"adhrit_sr_verma@wiza.org","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:06.945+05:30","updated_at":"2021-05-31T03:50:06.945+05:30"},{"id":331,"name":"Agnivesh Verma DO","email":"verma_do_agnivesh@dubuque-kuhlman.biz","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:09.261+05:30","updated_at":"2021-05-31T03:50:09.261+05:30"},{"id":401,"name":"Chandraketu Verma","email":"chandraketu_verma@douglas.io","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:10.296+05:30","updated_at":"2021-05-31T03:50:10.296+05:30"},{"id":410,"name":"Bhushan Verma","email":"bhushan_verma@jaskolski.com","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:10.442+05:30","updated_at":"2021-05-31T03:50:10.442+05:30"},{"id":639,"name":"Chandrabhaga Verma","email":"chandrabhaga_verma@kiehn.io","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:14.416+05:30","updated_at":"2021-05-31T03:50:14.416+05:30"},{"id":688,"name":"Trisha Verma","email":"verma_trisha@bahringer.name","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:15.343+05:30","updated_at":"2021-05-31T03:50:15.343+05:30"},{"id":735,"name":"Suresh Verma","email":"suresh_verma@towne-steuber.biz","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:16.117+05:30","updated_at":"2021-05-31T03:50:16.117+05:30"},{"id":782,"name":"Jyotsana Verma","email":"verma_jyotsana@boyle.org","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:16.942+05:30","updated_at":"2021-05-31T03:50:16.942+05:30"},{"id":883,"name":"Devasree Verma","email":"devasree_verma@metz-mohr.net","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:18.868+05:30","updated_at":"2021-05-31T03:50:18.868+05:30"},{"id":894,"name":"Bhadran Verma","email":"bhadran_verma@reilly.co","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:19.025+05:30","updated_at":"2021-05-31T03:50:19.025+05:30"},{"id":914,"name":"Ashlesh Verma","email":"verma_ashlesh@ankunding.com","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:19.322+05:30","updated_at":"2021-05-31T03:50:19.322+05:30"},{"id":952,"name":"Archan Verma","email":"archan_verma@zulauf.biz","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:19.859+05:30","updated_at":"2021-05-31T03:50:19.859+05:30"},{"id":1137,"name":"Akshaj Verma","email":"akshaj_verma@mcclure.org","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:23.327+05:30","updated_at":"2021-05-31T03:50:23.327+05:30"},{"id":1265,"name":"Gemine Verma","email":"gemine_verma@heathcote-bogisich.io","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:25.327+05:30","updated_at":"2021-05-31T03:50:25.327+05:30"},{"id":1271,"name":"Bhagwanti Verma","email":"bhagwanti_verma@muller-kuhlman.net","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:25.436+05:30","updated_at":"2021-05-31T03:50:25.436+05:30"},{"id":1285,"name":"Dinesh Verma","email":"dinesh_verma@tillman.com","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:25.716+05:30","updated_at":"2021-05-31T03:50:25.716+05:30"}]}`

let res = exp.evalExpression(md, menu, data);
console.log(res)
