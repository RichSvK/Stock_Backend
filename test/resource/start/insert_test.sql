INSERT INTO `broker_underwriter` (`broker_code`,`broker_name`)
VALUES ('FZ','Waterfront Sekuritas Indonesia')
    ,('HD','KGI Sekuritas Indonesia')
    ,('PG','Panca Global Sekuritas')
    ,('SH','Artha Sekuritas Indonesia');

INSERT INTO `stock_ipo` (`stock_code`,`price`,`ipo_shares`,`listed_shares`,`equity`,`warrant`,`nominal`,`mcb`,`is_affiliated`,`is_acceleration`,`is_new`,`is_full_commitment`,`is_not_involved_case`,`lock_up`,`subscribed_stock`)
VALUES ('STRK',100,1180000000,10721709000,47628929487,3245000000,12,0,0,0,1,1,1,1,31351707700);

INSERT INTO `ipo_detail` (`stock_code`,`uw_code`,`uw_shares`)
VALUES ('STRK','FZ',120360000)
    ,('STRK','HD',319780000)
    ,('STRK','PG',120360000)
    ,('STRK','SH',619500000);

