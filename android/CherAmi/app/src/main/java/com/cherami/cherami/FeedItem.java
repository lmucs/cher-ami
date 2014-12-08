package com.cherami.cherami;

import org.json.JSONObject;

public class FeedItem {
    public JSONObject msg;
    public FeedItem(){
        super();
    }

    public FeedItem(JSONObject msg) {
        super();
        this.msg = msg;
    }

}