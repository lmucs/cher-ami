package com.cherami.cherami;

import org.json.JSONObject;

public class FeedItem {
    public JSONObject msg;
    public String img;
    public FeedItem(){
        super();
    }

    public FeedItem(JSONObject msg) {
        super();
        this.msg = msg;
        this.img = "";
    }

    public FeedItem(JSONObject msg, String img) {
        super();
        this.msg = msg;
        this.img = img;
    }
}