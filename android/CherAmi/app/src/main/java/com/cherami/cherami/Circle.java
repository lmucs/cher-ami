package com.cherami.cherami;

import org.json.JSONException;
import org.json.JSONObject;

public class Circle {
    public JSONObject circle;

    public Circle(JSONObject circle) {
        this.circle = circle;
    }
    public JSONObject getCircle(){
        return this.circle;
    }
    public String getVisibility(){
        try {
            return this.circle.getString("visibility");
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return null;
    }

}
