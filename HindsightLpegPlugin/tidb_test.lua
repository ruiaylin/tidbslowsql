-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this
-- file, You can obtain one at http://mozilla.org/MPL/2.0/.
require "string"
local mysql = require "lpeg.mysql"
local tidb = require "lpeg.tidb"

function print_r ( t )
    local print_r_cache={}
    local function sub_print_r(t,indent)
        if (print_r_cache[tostring(t)]) then
            print(indent.."*"..tostring(t))
        else
            print_r_cache[tostring(t)]=true
            if (type(t)=="table") then
                for pos,val in pairs(t) do
                    if (type(val)=="table") then
                        print(indent.."["..pos.."] => "..tostring(t).." {")
                        sub_print_r(val,indent..string.rep(" ",string.len(pos)+8))
                        print(indent..string.rep(" ",string.len(pos)+6).."}")
                    elseif (type(val)=="string") then
                        print(indent.."["..pos..'] => "'..val..'"')
                    else
                        print(indent.."["..pos.."] => "..tostring(val))
                    end
                end
            else
                print(indent..tostring(t))
            end
        end
    end
    if (type(t)=="table") then
        print(tostring(t).." {")
        sub_print_r(t,"  ")
        print("}")
    else
        sub_print_r(t,"  ")
    end
    print()
end

local tidb_slow_query_log = {
[[
# Time: 2019-09-22T22:32:12.080095163+08:00
# Txn_start_ts: 411346595049111687
# User: reply_wr@10.27.75.75
# Conn_ID: 700288
# Query_time: 0.686471601
]],
[[
# Time: 2019-09-22T22:32:12.111440145+08:00
# Txn_start_ts: 411346594852503760
# User: reply_wr@10.27.75.72
# Conn_ID: 700288
# Query_time: 0.686471601
# Process_time: 0.002 Wait_time: 0.68 Request_count: 6 Total_keys: 32 Process_keys: 30
# DB: reply
# Index_ids: [1]
# Is_internal: false
# Digest: fb0a8169a2b6dfd2c455d26335ef4238986bbd877db0cdd50b0c6c80dc9dc11f
# Stats: ClubReply:411346498422833348
# Num_cop_tasks: 2
# Cop_proc_avg: 0.0025 Cop_proc_p90: 0.005 Cop_proc_max: 0.005
# Cop_wait_avg: 0.335 Cop_wait_p90: 0.67 Cop_wait_max: 0.67
# Mem_max: 9740
SELECT  replyid,rContent as 'r_content',rReplyDate as 'replydate',rLastEditDate as 'r_lasteditdate',  rMemberId as 'r_memberid' ,rMemberName as 'r_membername',rTargetReplyid as 'targetreplyid',  rFloor as 'floor',rdelete,rishtml,risuserhtml,rSource as 'rsource'  FROM ClubReply r   WHERE rTopicId=78308330  ORDER BY replyid limit 20;
]],
[[
# Time: 2019-09-22T22:32:12.080095163+08:00
# Txn_start_ts: 411346595049111687
# User: reply_wr@10.27.75.75
# Conn_ID: 700288
# Query_time: 0.686471601
# Process_time: 0.002 Wait_time: 0.68 Request_count: 6 Total_keys: 32 Process_keys: 30
# DB: reply
]],[[
# Time: 2019-09-22T22:32:12.094247985+08:00
# Txn_start_ts: 411346595049111655
# User: reply_wr@10.27.74.47
# Conn_ID: 557475
# Query_time: 0.709839379
# Wait_time: 0.71 Request_count: 1 Total_keys: 18 Process_keys: 17
# DB: reply
# Index_ids: [1]
# Is_internal: false
# Digest: 71a1128b093ca91a92121beb80c52cf38bcdeec345667b3d05c562832615bc6a
# Stats: ClubReply:411346498422833348
# Num_cop_tasks: 1
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0.71 Cop_wait_p90: 0.71 Cop_wait_max: 0.71
# Mem_max: 440
SELECT COUNT(replyId) FROM ClubReply WHERE rTopicId=62039947;
]],
[[
# Time: 2019-09-22T22:32:12.094247985+08:00
# Txn_start_ts: 411346595049111655
# User: reply_wr@10.27.74.47
# Conn_ID: 557475
# Query_time: 0.709839379
# Wait_time: 0.71 Request_count: 1 Total_keys: 18 Process_keys: 17
# DB: reply
# Index_ids: [1]
# Is_internal: false
# Digest: 71a1128b093ca91a92121beb80c52cf38bcdeec345667b3d05c562832615bc6a
# Stats: ClubReply:411346498422833348
# Num_cop_tasks: 1
# Cop_proc_avg: 0.1 Cop_proc_p90: 0.1 Cop_proc_max: 0.1
# Cop_wait_avg: 0.71 Cop_wait_p90: 0.71 Cop_wait_max: 0.71
# Mem_max: 440
SELECT COUNT(replyId) FROM ClubReply WHERE rTopicId=62039947;
]],[[
# Time: 2019-09-22T22:32:12.094247985+08:00
# Txn_start_ts: 411346595049111655
# User: reply_wr@10.27.74.47
# Conn_ID: 557475
# Query_time: 0.709839379
# Wait_time: 0.71 Request_count: 1 Total_keys: 18 Process_keys: 17
# DB: reply
# Index_ids: [1]
# Is_internal: false
# Digest: 71a1128b093ca91a92121beb80c52cf38bcdeec345667b3d05c562832615bc6a
# Stats: ClubReply:411346498422833348
# Num_cop_tasks: 1
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0.71 Cop_wait_p90: 0.71 Cop_wait_max: 0.71
# Mem_max: 440
SELECT COUNT(replyId) FROM ClubReply WHERE rTopicId=62039947;
]],[[
# Time: 2019-09-24T08:57:48.639983212+08:00
# Txn_start_ts: 411379084285181980
# User: reply_wr@10.27.75.71
# Conn_ID: 814993
# Query_time: 0.663635661
# Process_time: 0.663 Request_count: 1 Total_keys: 769302 Process_keys: 769301
# DB: reply
# Index_ids: [5]
# Is_internal: false
# Digest: a003fdd85bda18195e505ccc50c3d9e9c95e09a9375a5c8b427bfb4b35296755
# Stats: ClubReply:411379053339082818
# Num_cop_tasks: 1
# Cop_proc_avg: 0.663 Cop_proc_p90: 0.663 Cop_proc_max: 0.663
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
# Mem_max: 141
SELECT count(id) FROM ClubReply r  WHERE rTopicId=19582947;
]],[[
# Time: 2019-09-22T22:32:19.913494105+08:00
# Txn_start_ts: 411346597159370780
# User: autopricewriter@10.27.3.36
# Conn_ID: 657968
# Query_time: 0.49688529
# Wait_time: 0.496 Request_count: 3 Total_keys: 4 Process_keys: 2
# DB: autoprice
# Index_ids: [10]
# Is_internal: false
# Digest: e9f8139213cd596c1fc994ddae0496b5518f644bb7b6c92f16437cb56e29db44
# Stats: car_price:411346506287153225,car_installment:411343191515594801
# Num_cop_tasks: 3
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0.165333333 Cop_wait_p90: 0.495 Cop_wait_max: 0.495
# Mem_max: 83565
 SELECT  IFNULL(b.poundage,0) poundage,IFNULL(b.down_payment,0)down_payment,IFNULL(b.interest,0) intere
st,IFNULL(b.staging_number,0)staging_number, a.id, a.member_id, a.member_name, a.brand_id, a.series_id,
 a.spec_id, a.naked_price, a.purchase_tax, a.insurer, a.use_tax, a.traffic_insurance, a.license_fee, a.
full_price, a.sales_pack, date_format(a.shopping_time,'%Y-%m-%d') shopping_time, a.province_id, a.city_
id, a.county_id, a.dealer_id, a.feel_level, a.feel_content, a.invoice_url, a.status, a.ip, a.client_typ
e, a.is_hide, a.price_type, date_format(a.gmt_create,'%Y-%m-%d %H:%I:%S') gmt_create, date_format(a.gmt
_modify,'%Y-%m-%d %H:%I:%S') gmt_modify, a.pk_uuid, a.is_reward FROM car_price a LEFT JOIN car_installm
ent b ON a.id = b.car_price_id where 1=1  and a.is_hide = 0 and a.spec_id=38323  and a.city_id=441900
and a.`status` = 1 ORDER BY   a.shopping_time DESC limit 0 ,10;
]],[[
# Time: 2019-09-24T07:30:24.13313912+08:00
# Txn_start_ts: 411377709484081242
# User: reply_wr@10.27.75.118
# Conn_ID: 866317
# Query_time: 0.586036754
# Wait_time: 0.001 Backoff_time: 0.582 Request_count: 2 Total_keys: 3 Process_keys: 2
# DB: Reply
# Index_ids: [3]
# Is_internal: false
# Digest: 6b35d2667b56862d903e5a86cfaf2c21b56c0ca884ed7bd0e946d498fba5efd3
# Stats: ClubReply:411377641012068450
# Num_cop_tasks: 2
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0.0005 Cop_wait_p90: 0.001 Cop_wait_max: 0.001
# Mem_max: 4328
SELECT rTopicId AS topicid, replyId AS topicid1, rMemberId AS tmemberid, rMemberName AS tMemberName,         '' AS ttitle, rContent AS autopicContent, rReplyDate AS tdate, 0 AS tbbsid, '' AS bbsName,         rFloor AS `floor`, '' AS topicUrl, rIP AS topicIp, 0 AS replyCount, 0 AS topicType,         0 AS IfAccept,         CASE WHEN CreateType = 0 THEN 0 WHEN CreateType = 600 THEN 3 ELSE 2 END AS replySource,         0 AS statusCode, CURRENT_TIMESTAMP( ) AS statusTime  FROM ClubReply  WHERE replyId = 1256583658;
]],[[
# Time: 2019-09-24T08:10:43.71929167+08:00
# Txn_start_ts: 411378343846346760
# User: autopricewriter@10.27.3.27
# Conn_ID: 857785
# Query_time: 0.306906354
# Process_time: 0.516 Wait_time: 0.006 Backoff_time: 0.1 Request_count: 31 Total_keys: 29640 Process_keys: 29010
# DB: autoprice
# Index_ids: [22]
# Is_internal: false
# Digest: 84844ec2d649dae6f674d270c9413235b4882f3b6c272703415fc9760fd4bc90
# Stats: car_price:411378278061310071
# Num_cop_tasks: 31
# Cop_proc_avg: 0.016645161 Cop_proc_p90: 0.041 Cop_proc_max: 0.076
# Cop_wait_avg: 0.000193548 Cop_wait_p90: 0.001 Cop_wait_max: 0.001
# Mem_max: 15772868
select carprice0_.id as id1_16_, carprice0_.brand_id as brand_id2_16_, carprice0_.city_id as city_id3_16_, carprice0_.client_type as client_t4_16_, carprice0_.county_id as county_i5_16_, carprice0_.dealer_id as dealer_i6_16_, carprice0_.feel_content as feel_con7_16_, carprice0_.feel_level as feel_lev8_16_, carprice0_.full_price as full_pri9_16_, carprice0_.gmt_create as gmt_cre10_16_, carprice0_.gmt_modify as gmt_mod11_16_, carprice0_.insurer as insurer12_16_, carprice0_.invoice_url as invoice13_16_, carprice0_.ip as ip14_16_, carprice0_.is_hide as is_hide15_16_, carprice0_.is_reward as is_rewa16_16_, carprice0_.license_fee as license17_16_, carprice0_.member_id as member_18_16_, carprice0_.member_name as member_19_16_, carprice0_.naked_price as naked_p20_16_, carprice0_.pk_uuid as pk_uuid21_16_, carprice0_.price_type as price_t22_16_, carprice0_.province_id as provinc23_16_, carprice0_.purchase_tax as purchas24_16_, carprice0_.remark1 as remark25_16_, carprice0_.sales_pack as sales_p26_16_, carprice0_.series_id as series_27_16_, carprice0_.shopping_time as shoppin28_16_, carprice0_.spec_id as spec_id29_16_, carprice0_.status as status30_16_, carprice0_.traffic_insurance as traffic31_16_, carprice0_.use_tax as use_tax32_16_ from car_price carprice0_ where carprice0_.series_id=2922 and carprice0_.status=1 and carprice0_.is_hide=0 order by carprice0_.shopping_time desc limit 9800, 20;
]],[[
# Time: 2019-09-23T23:23:29.924855191+08:00
# Txn_start_ts: 411370051156836382
# User: autopricewriter@10.27.3.30
# Conn_ID: 831973
# Query_time: 0.606404605
# Process_time: 0.003 Wait_time: 0.602 Request_count: 15 Total_keys: 40 Process_keys: 26
# DB: autoprice
# Index_ids: [35]
# Is_internal: false
# Digest: e9f8139213cd596c1fc994ddae0496b5518f644bb7b6c92f16437cb56e29db44
# Stats: car_price:411369941921431601,car_installment:411368518191153220
# Num_cop_tasks: 15
# Cop_proc_avg: 0.0002 Cop_proc_p90: 0.001 Cop_proc_max: 0.001
# Cop_wait_avg: 0.040133333 Cop_wait_p90: 0 Cop_wait_max: 0.602
# Mem_max: 64292
 SELECT  IFNULL(b.poundage,0) poundage,IFNULL(b.down_payment,0)down_payment,IFNULL(b.interest,0) interest,IFNULL(b.staging_number,0)staging_number, a.id, a.member_id, a.member_name, a.brand_id, a.series_id, a.spec_id, a.naked_price, a.purchase_tax, a.insurer, a.use_tax, a.traffic_insurance, a.license_fee, a.full_price, a.sales_pack, date_format(a.shopping_time,'%Y-%m-%d') shopping_time, a.province_id, a.city_id, a.county_id, a.dealer_id, a.feel_level, a.feel_content, a.invoice_url, a.status, a.ip, a.client_type, a.is_hide, a.price_type, date_format(a.gmt_create,'%Y-%m-%d %H:%I:%S') gmt_create, date_format(a.gmt_modify,'%Y-%m-%d %H:%I:%S') gmt_modify, a.pk_uuid, a.is_reward FROM car_price a LEFT JOIN car_installment b ON a.id = b.car_price_id where 1=1  and a.is_hide = 0 and a.spec_id=29124  and a.city_id=510100  and a.`status` = 1 ORDER BY   a.shopping_time DESC limit 0 ,1000;
]],
[[
# Time: 2019-09-22T23:12:27.940108183+08:00
# Txn_start_ts: 411347228323479557
# User: reply_wr@10.27.74.25
# Conn_ID: 451702
# Query_time: 0.388001568
# Wait_time: 0.823 Request_count: 3 Total_keys: 202 Process_keys: 201 Prewrite_time: 0.000825655 Commit_time: 0.000666926 Get_commit_ts_time: 7.2763e-05 Write_keys: 1 Write_size: 230 Prewrite_region: 1
# DB: reply
# Index_ids: [1,1]
# Is_internal: false
# Digest: 76c0b053aaec74395e6b48b133f38658cf0fb7344eb983609c5d9e067c77e60b
# Stats: ClubReply:411347124095025163
# Num_cop_tasks: 3
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0.274333333 Cop_wait_p90: 0.438 Cop_wait_max: 0.438
# Mem_max: 13280
                UPDATE  ClubReply SET rFloor = (SELECT COUNT(id) FROM  ClubReply WHERE rTopicId = 82630975 AND replyId <= 1256312528) WHERE rTopicId = 82630975 AND replyId = 1256312528;
]],[[
# Time: 2019-09-22T22:34:39.139163261+08:00
# Txn_start_ts: 411346633269706795
# User: reply_wr@10.27.74.29
# Conn_ID: 704845
# Query_time: 1.969641008
# Process_time: 0.002 Wait_time: 1.967 Request_count: 4 Total_keys: 122 Process_keys: 120
# DB: Reply
# Index_ids: [1,1]
# Is_internal: false
# Digest: 5d0e10bd654adc91e6185b55e1f8888e94cbcaa950d351f8821261b76d17279a
# Stats: ClubReply:411346498422833348
# Num_cop_tasks: 4
# Cop_proc_avg: 0.0005 Cop_proc_p90: 0.001 Cop_proc_max: 0.001
# Cop_wait_avg: 0.49175 Cop_wait_p90: 1.967 Cop_wait_max: 1.967
# Mem_max: 17012
SELECT replyId, rTopicId, rContent, rMemberId, rMemberName, rReplyDate, rIP, rDelete, rLastEditDate, rIsHTML, rFloor, rIsUserHTML, CreateType, 
rTarget, rTargetReplyid FROM ClubReply WHERE rTopicId =82360066 AND replyId > ( SELECT MAX(replyId) FROM ( SELECT replyId FROM 
    ClubReply WHERE rTopicId =82360066 ORDER BY replyId LIMIT 80 ) AS K ) ORDER BY replyId  LIMIT 20;
]],[[
# Time: 2019-09-23T23:01:52.403565526+08:00
# Txn_start_ts: 411369710998781988
# User: reply_wr@10.27.75.71
# Conn_ID: 836950
# Query_time: 0.689265905
# Process_time: 0.692 Wait_time: 0.001 Request_count: 14 Total_keys: 118 Process_keys: 116
# DB: reply
# Index_ids: [7]
# Is_internal: false
# Digest: 886bfa770c8aa560e2e9da2d1313deace8fcdae80eea13d9f213e413925f9a55
# Stats: ClubReply:411369619576062050
# Num_cop_tasks: 14
# Cop_proc_avg: 0.049428571 Cop_proc_p90: 0.001 Cop_proc_max: 0.687
# Cop_wait_avg: 7.1428e-05 Cop_wait_p90: 0 Cop_wait_max: 0.001
# Mem_max: 5728
SELECT rFloor as 'floor'  FROM ClubReply r   WHERE rTopicId=75448862 AND rMemberId=1978466  ORDER BY replyid;
]],[[
# Time: 2019-09-23T06:05:22.715461907+08:00
# Txn_start_ts: 411353722914865184
# User: reply_wr@10.27.74.43
# Conn_ID: 708420
# Query_time: 0.665928559
# Wait_time: 0.665 Request_count: 1 Total_keys: 1
# DB: reply
# Index_ids: [1]
# Is_internal: false
# Digest: e5a0e78325db64bc492986eadfae8eeedcf68b204b37a967eac961be828a1c46
# Stats: ClubReply:411353576284618775
# Num_cop_tasks: 1
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0.665 Cop_wait_p90: 0.665 Cop_wait_max: 0.665
# Mem_max: 125
SELECT replyid,rContent as 'r_content',rReplyDate as 'replydate',rLastEditDate as 'r_lasteditdate',rMemberId as 'r_memberid' ,rMemberName as 'r_membername',rTargetReplyid as 'targetreplyid', rFloor as 'floor',rdelete,rishtml,risuserhtml,rSource as 'rsource' FROM ClubReply r WHERE rTopicId=83794488 ORDER BY replyid desc limit 1000;
]],[[
# Time: 2019-08-14T09:26:59.487776265+08:00
# Txn_start_ts: 410450924122144769
# User: root@127.0.0.1
# Conn_ID: 3086
# Query_time: 1.527627037
# Process_time: 0.07 Request_count: 1 Total_keys: 131073 Process_keys: 131072 Prewrite_time: 0.335415029 Commit_time: 0.032175429 Get_commit_ts_time: 0.000177098 Local_latch_wait_time: 0.106869448 Write_keys: 131072 Write_size: 3538944 Prewrite_region: 1
# DB: test_1
# Is_internal: false
# Digest: 50a2e32d2abbd6c1764b1b7f2058d428ef2712b029282b776beb9506a365c0f1
# Stats: t:pseudo
# Num_cop_tasks: 1
# Cop_proc_avg: 0.07 Cop_proc_p90: 0.07 Cop_proc_max: 0.07 Cop_proc_addr: 172.16.5.87:20171
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr: 172.16.5.87:20171
# Mem_max: 525211
# Succ: false
insert into t select * from t;
]],[[
# Time: 2019-09-22T23:27:59.518031776+08:00
# Txn_start_ts: 0
# Query_time: 28.409603522
# Is_internal: true
# Digest: 2917a8c8bb18aff03c6dcb6b1afbbab3067e8d6934dc9f0817bf9ae95c543f1a
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
analyze table `autoprice`.`car_price_statistics`;
]],[[
# Time: 2019-09-23T21:48:57.889272146+08:00
# Txn_start_ts: 0
# User: dba@127.0.0.1
# Conn_ID: 831556
# Query_time: 2.211217815
# DB: information_schema
# Is_internal: false
# Digest: 9d8d1c27d43986bbe6bc047ca8b48e777b082b434033f32c3d9431740b31e35e
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
            select Time,                 Txn_start_ts,                 User,                 Conn_ID,                 Query_time,                 Process_time,                 Wait_time,                 Backoff_time,                 Request_count,                 Total_keys,                 Process_keys,                 DB,                 Index_ids,                 Is_internal,                 Digest,                 Stats,                 Cop_proc_avg,                 Cop_proc_p90,                 Cop_proc_max,                 Cop_wait_avg,                 Cop_wait_p90,                 Cop_wait_max,                 Mem_max,                 Query               from slow_query              where time>date_sub(now(),interval 10 minute)                  and query not like '%slow_query%';
]],[[
# Time: 2019-09-22T22:51:20.44982795+08:00
# Txn_start_ts: 411346896187031569
# User: autopricewriter@10.27.2.32
# Conn_ID: 657128
# Query_time: 0.312873027
# Prewrite_time: 0.310802323 Commit_time: 0.001882215 Get_commit_ts_time: 8.9563e-05 Write_keys: 81 Write_size: 3829 Prewrite_region: 7
# DB: autoprice
# Is_internal: false
# Digest: 9505cacb7c710ed17125fcc6cb3669e8ddca6c8cd8af6a31f6b3cd64604c3098
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
commit;
]],[[
# Time: 2019-09-23T10:41:24.261911464+08:00
# Txn_start_ts: 411358064491364387
# User: reply_wr@10.27.74.26
# Conn_ID: 757559
# Query_time: 0.425277471
# Request_count: 1 Total_keys: 10 Process_keys: 9
# DB: Reply
# Index_ids: [1]
# Is_internal: false
# Digest: 76d4c8de80230dc532e28a55e193ddb0233be15dfb7d823de214ec36963e1e93
# Stats: ClubReply:411357980198961183
# Num_cop_tasks: 1
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
# Mem_max: 209
SELECT COUNT(1)  FROM  ClubReply WHERE rTopicId=82694806;
]],[[
# Time: 2019-09-23T11:13:10.432803591+08:00
# Txn_start_ts: 0
# User: dba@10.27.74.132
# Conn_ID: 780205
# Query_time: 1.7672986499999999
# DB: information_schema
# Is_internal: false
# Digest: b10085fa74b7b7012a8ca27ad260b7033b0284633a54c99747e45c68c3d16c5b
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
select count(*) from slow_query where time>'2019-09-23 9:30:00';
]],[[
# Time: 2019-09-23T11:13:47.67622078+08:00
# Txn_start_ts: 0
# User: dba@10.27.74.132
# Conn_ID: 780205
# Query_time: 2.084945073
# DB: information_schema
# Is_internal: false
# Digest: f1de62a09464e8a9f2c70b42ac315cbf8a6285222948e4b1b7b5dc02c77e492c
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
select *  from slow_query where time>'2019-09-23 9:30:00' limit 10;
]],[[
# Time: 2019-09-02T16:21:10.385337925+08:00
# Txn_start_ts: 410887770574946305
# Query_time: 15.483088376
# Prewrite_time: 11.174731535 Commit_time: 1.875035557 Get_commit_ts_time: 0.000378074 Local_latch_wait_time: 1.5926e-05 Write_keys: 16 Write_size: 453491 Prewrite_region: 1 Txn_retry: 1
# Is_internal: true
# Digest: 8cea8954b5a9d7e444feb61c1f5c240c0088d80c8cd4dfbb8a75a36a9579c799
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0 Cop_proc_addr:
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr:
# Succ: false
insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (81, 0, 31, 0, 12654317),(81, 0, 4, 0, 12977),(81, 0, 7, 0, 43162),(81, 0, 11, 0, 4073),(81, 0, 2, 0, 39566),(81, 0, 30, 0, 432031),(81, 0, 23, 0, 9320),(81, 0, 10, 0, 4115),(81, 0, 3, 0, 11386),(81, 0, 29, 0, 39555),(81, 0, 5, 0, 10529),(81, 0, 8, 0, 47520980),(81, 0, 6, 0, 17708),(81, 0, 9, 0, 12716454),(81, 0, 24, 0, 4660),(81, 0, 22, 0, 590) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-08-17T21:04:54.922748376+08:00
# Txn_start_ts: 410529848848285701
# Query_time: 3.036920735
# Prewrite_time: 0.058361101 Commit_time: 0.103531395 Get_commit_ts_time: 0.000149317 Local_latch_wait_time: 3.1903e-05 Write_keys: 10 Write_size: 308593 Prewrite_region: 2
# Is_internal: true
# Digest: 85aabb5a0425977cd016e15f3de5bb40e3bba1628c178ba32e386fefd4f42fa1
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0 Cop_proc_addr:
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr:
# Succ: false
insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 1, 0, 12416416),(237, 0, 15, 0, 34145144),(237, 0, 2, 0, 1552052),(237, 0, 9, 0, 4500292),(237, 0, 6, 0, 1940065),(237, 0, 3, 0, 2328078),(237, 0, 4, 0, 37637257),(237, 0, 16, 0, 12416416),(237, 0, 7, 0, 1940065),(237, 0, 5, 0, 1940065) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-08-17T21:04:02.863737213+08:00
# Txn_start_ts: 410529835911479297
# Query_time: 0.344850466
# Prewrite_time: 3.424110363 Commit_time: 0.000522605 Get_commit_ts_time: 0.0002014 Resolve_lock_time: 0.153877379 Local_latch_wait_time: 3.1218e-05 Write_keys: 10 Write_size: 308593 Prewrite_region: 2 Txn_retry: 2
# Is_internal: true
# Digest: 85aabb5a0425977cd016e15f3de5bb40e3bba1628c178ba32e386fefd4f42fa1
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0 Cop_proc_addr:
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr:
# Succ: false
insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 15, 0, 30537232),(237, 0, 2, 0, 1388056),(237, 0, 1, 0, 11104448),(237, 0, 9, 0, 4025528),(237, 0, 4, 0, 33660356),(237, 0, 16, 0, 11104448),(237, 0, 3, 0, 2082084),(237, 0, 6, 0, 1735070),(237, 0, 7, 0, 1735070),(237, 0, 5, 0, 1735070) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-10-23T14:20:01.689536516+08:00
# Txn_start_ts: 412040980470169608
# User: gzconcern_wr@10.27.3.28
# Conn_ID: 408366
# Query_time: 0.10289279
# Local_latch_wait_time: 0.002043375 Write_keys: 4 Write_size: 178 Txn_retry: 1
# DB: gzconcern
# Is_internal: false
# Digest: 33310ef1dc065104278e13aa8d9da0a80f334b22f72026b378b818bfb10feb9d
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0 Cop_proc_addr: 
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr: 
# Succ: true
insert into user_follow_relation ( original_uid, target_uid, created_stime, modified_stime ) values ( 37306023, 71109611, '2019-10-23 14:20:02', '2019-10-23 14:20:01.62' ) ON DUPLICATE KEY UPDATE modified_stime = CURRENT_TIMESTAMP;
]]
}

local function tidb_test1()
    print(tidb_slow_query_log[1])
    local t = tidb.tidb_slow_sql:match(tidb_slow_query_log[1])
    if not t then return error("no match") end
    print("data:", t)
    print_r( t )
end

local function tidb_test2(i)
    print(tidb_slow_query_log[i])
    local t = tidb.tidb21_slow_sql_grammar:match(tidb_slow_query_log[i])
    if not t then return error("no match") end
    print("data:", t)
    print_r( t )
end

--tidb_test1()
for i = 1 , #tidb_slow_query_log do
   print("---------------------------------------------------------")
   tidb_test2(i)
end



