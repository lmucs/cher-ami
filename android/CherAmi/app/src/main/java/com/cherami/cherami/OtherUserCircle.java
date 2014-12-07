package com.cherami.cherami;

import org.json.JSONObject;

/**
 * Created by Geoff on 11/30/2014.
 */
public class OtherUserCircle {
    public JSONObject circle;

    public OtherUserCircle(JSONObject circle) {
        this.circle = circle;
    }

    public JSONObject getCircle(){
        return this.circle;
    }
}
