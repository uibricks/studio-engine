# Steps to execute the JS file

1. Start Node Repl

2. Import JS File
    ```
    exp=require("./expressions_utility.js")
    ```

3. Set Metadata
    ```
    md={
      "tYu6c8jd1BHprNcB4mztc": {
        "id": "tYu6c8jd1BHprNcB4mztc",
        "name": "GetAllNames",
        "raw": "[phd8JZrz5Z95uT95-Lmx7]",
        "refs": {
          "phd8JZrz5Z95uT95-Lmx7": {
            "path": ["data", "#", "name"],
            "type": "stringArray"
          }
        },
        "type": "stringArray"
      }
    }
    ```
 
4. Set Menu
    ```
    menu=[{"id": "tYu6c8jd1BHprNcB4mztc"}]
    ```

5. Set Data
    ```
    data=`{"code":200,"meta":{"pagination":{"total":19,"pages":1,"page":1,"limit":20}},"data":[{"id":8,"name":"Gitanjali Verma","email":"gitanjali_verma@hayes.com","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:04.174+05:30","updated_at":"2021-05-31T03:50:04.174+05:30"},{"id":70,"name":"Deependra Verma","email":"verma_deependra@marks-fahey.name","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:05.220+05:30","updated_at":"2021-05-31T03:50:05.220+05:30"},{"id":97,"name":"Tapan Verma","email":"tapan_verma@bernhard.net","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:05.618+05:30","updated_at":"2021-05-31T03:50:05.618+05:30"},{"id":178,"name":"Adhrit Verma Sr.","email":"adhrit_sr_verma@wiza.org","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:06.945+05:30","updated_at":"2021-05-31T03:50:06.945+05:30"},{"id":331,"name":"Agnivesh Verma DO","email":"verma_do_agnivesh@dubuque-kuhlman.biz","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:09.261+05:30","updated_at":"2021-05-31T03:50:09.261+05:30"},{"id":401,"name":"Chandraketu Verma","email":"chandraketu_verma@douglas.io","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:10.296+05:30","updated_at":"2021-05-31T03:50:10.296+05:30"},{"id":410,"name":"Bhushan Verma","email":"bhushan_verma@jaskolski.com","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:10.442+05:30","updated_at":"2021-05-31T03:50:10.442+05:30"},{"id":639,"name":"Chandrabhaga Verma","email":"chandrabhaga_verma@kiehn.io","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:14.416+05:30","updated_at":"2021-05-31T03:50:14.416+05:30"},{"id":688,"name":"Trisha Verma","email":"verma_trisha@bahringer.name","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:15.343+05:30","updated_at":"2021-05-31T03:50:15.343+05:30"},{"id":735,"name":"Suresh Verma","email":"suresh_verma@towne-steuber.biz","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:16.117+05:30","updated_at":"2021-05-31T03:50:16.117+05:30"},{"id":782,"name":"Jyotsana Verma","email":"verma_jyotsana@boyle.org","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:16.942+05:30","updated_at":"2021-05-31T03:50:16.942+05:30"},{"id":883,"name":"Devasree Verma","email":"devasree_verma@metz-mohr.net","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:18.868+05:30","updated_at":"2021-05-31T03:50:18.868+05:30"},{"id":894,"name":"Bhadran Verma","email":"bhadran_verma@reilly.co","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:19.025+05:30","updated_at":"2021-05-31T03:50:19.025+05:30"},{"id":914,"name":"Ashlesh Verma","email":"verma_ashlesh@ankunding.com","gender":"Female","status":"Inactive","created_at":"2021-05-31T03:50:19.322+05:30","updated_at":"2021-05-31T03:50:19.322+05:30"},{"id":952,"name":"Archan Verma","email":"archan_verma@zulauf.biz","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:19.859+05:30","updated_at":"2021-05-31T03:50:19.859+05:30"},{"id":1137,"name":"Akshaj Verma","email":"akshaj_verma@mcclure.org","gender":"Male","status":"Active","created_at":"2021-05-31T03:50:23.327+05:30","updated_at":"2021-05-31T03:50:23.327+05:30"},{"id":1265,"name":"Gemine Verma","email":"gemine_verma@heathcote-bogisich.io","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:25.327+05:30","updated_at":"2021-05-31T03:50:25.327+05:30"},{"id":1271,"name":"Bhagwanti Verma","email":"bhagwanti_verma@muller-kuhlman.net","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:25.436+05:30","updated_at":"2021-05-31T03:50:25.436+05:30"},{"id":1285,"name":"Dinesh Verma","email":"dinesh_verma@tillman.com","gender":"Male","status":"Inactive","created_at":"2021-05-31T03:50:25.716+05:30","updated_at":"2021-05-31T03:50:25.716+05:30"}]}`
    ```

6. Invoke expression function
   1. To evaluate multiple expressions
    ```
    exp.evalExpressions(md,menu,data)
    ```
    use below to view nested objects
    ```
    console.log(util.inspect(exp.evalExpressions(md,menu,data), false, null, true))
    ```
   2. To evaluate single expressions [metadatas contain expression for the whole page/repo, while md is a single expression]
    ```
   exp.evalExpression(metadatas,md,data)
    ```
   use below to view nested objects
    ```
    console.log(util.inspect(exp.evalExpression(metadatas,md,data), false, null, true))
    ```
   

7. Result
    ```
    {
      "tYu6c8jd1BHprNcB4mztc": [
        "Gitanjali Verma",
        "Deependra Verma",
        "Tapan Verma",
        "Adhrit Verma Sr.",
        "Agnivesh Verma DO",
        "Chandraketu Verma",
        "Bhushan Verma",
        "Chandrabhaga Verma",
        "Trisha Verma",
        "Suresh Verma",
        "Jyotsana Verma",
        "Devasree Verma",
        "Bhadran Verma",
        "Ashlesh Verma",
        "Archan Verma",
        "Akshaj Verma",
        "Gemine Verma",
        "Bhagwanti Verma",
        "Dinesh Verma"
      ]
    }
    ```


# For more examples, go-to : 
    https://uibricks.atlassian.net/wiki/spaces/~990241242/pages/2543124520/Expressions

# Invoke 'evalExpression' from expression-editor

1. Data
   ```
   data = `{"billingInformation":{"autoPay":false,"billedToDate":"05/22/2020","lastPaymentDate":"04/13/2020","paymentAmount":45.28,"paymentMethod":"Wells Fargo Checking","paymentMode":"Monthly"},"deathBenefitAmount":1000000,"effectiveDate":"04/23/2010","faceAmount":500000,"health":"Good","maturityDate":"04/22/2030","ownerInformation":{"corporateOwnerName":"Watchtower Corporate","id":"56789","individualOwner":{"address":{"city":"Portland","country":"USA","state":"OR","street1":"4430 Heron Way","street2":null,"zip4":null,"zip5":"97204"},"name":{"first":"Billy","last":"Daugherty","middle":"P","prefix":"Mr.","suffix":null},"policyRelationshipCode":""},"partyId":""},"payorInformation":{"corporatePayorName":"Watchtower Corporate","individualPayor":{"address":{"city":"Portland","country":"USA","state":"OR","street2":null,"zip4":null,"zip5":"97204"},"id":"5322","name":{"first":"Billy","last":"Daugherty","middle":"P","prefix":"Mr.","suffix":null},"policyRelationshipCode":"Payor"}},"planName":"Term Life Insurance","policyId":"456-346-4715","premiumDueDate":"05/23/2020","productType":"Life Insurance","status":"Active","termLength":"20 years"}`
   ```
2. Saved Metdata
   ```
   metadatas=
   {
      "complete-address-id": {
        "id": "complete-address-id",
        "name": "completeAddress",
        "raw": "ConcatDelimiter(',',[aRZY-8BHxI7cN6ju6Jk1u],[gIS3SwPUVNBB2kfpv7xuJ],[eH-QMUdW_9L9kOGRN1BP7],[XLOHP1ZyDWYJ9qeqwK_8u])",
        "refs": {
          "XLOHP1ZyDWYJ9qeqwK_8u": {
            "path": [
              "ownerInformation",
              "individualOwner",
              "address",
              "zip5"
            ],
            "type": "string"
          },
          "aRZY-8BHxI7cN6ju6Jk1u": {
            "path": [
              "ownerInformation",
              "individualOwner",
              "address",
              "street1"
            ],
            "type": "string"
          },
          "eH-QMUdW_9L9kOGRN1BP7": {
            "path": [
              "ownerInformation",
              "individualOwner",
              "address",
              "state"
            ],
            "type": "string"
          },
          "gIS3SwPUVNBB2kfpv7xuJ": {
            "path": [
              "ownerInformation",
              "individualOwner",
              "address",
              "city"
            ],
            "type": "string"
          }
        },
        "type": "string"
      }
   }
   ```
   
3. To-Evaluate Metadata
   ```
   md=
   {
     "owner-name-id": {
       "id": "owner-name-id",
       "name": "fullName",
       "nestedRefs": [
         "complete-address-id"
       ],
       "raw": "ConcatDelimiter(' lives at ', ConcatDelimiter(' ',[XLOHP1ZyDWYJ9qeqwK_8u2], [aRZY-8BHxI7cN6ju6Jk1u2],[eH-QMUdW_9L9kOGRN1BP72],[gIS3SwPUVNBB2kfpv7xuJ2]), [complete-address-id])",
       "refs": {
         "XLOHP1ZyDWYJ9qeqwK_8u2": {
           "path": [
             "ownerInformation",
             "individualOwner",
             "name",
             "prefix"
           ],
           "type": "string"
         },
         "aRZY-8BHxI7cN6ju6Jk1u2": {
           "path": [
             "ownerInformation",
             "individualOwner",
             "name",
             "first"
           ],
           "type": "string"
         },
         "eH-QMUdW_9L9kOGRN1BP72": {
           "path": [
             "ownerInformation",
             "individualOwner",
             "name",
             "middle"
           ],
           "type": "string"
         },
         "gIS3SwPUVNBB2kfpv7xuJ2": {
           "path": [
             "ownerInformation",
             "individualOwner",
             "name",
             "last"
           ],
           "type": "string"
         }
       },
       "type": "string"
     }
   }
   ```
4. Invoke 'evalExpression'
   ```
   exp.evalExpression(metadatas,md,data)
   ```
5. Result
   ```
   { 'owner-name-id':
     'Mr. Billy P Daugherty lives at 4430 Heron Way,Portland,OR,97204' }
   ```