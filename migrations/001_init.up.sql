CREATE TABLE stats
(
    timestamp  DateTime64,
    useragent  String,
    ip_address String
) Engine = Memory;
