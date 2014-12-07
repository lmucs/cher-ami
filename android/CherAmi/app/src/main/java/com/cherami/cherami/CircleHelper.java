package com.cherami.cherami;

import android.app.Activity;
import android.content.res.AssetManager;
import android.content.res.Resources;
import android.util.Log;
import android.view.View;
import android.widget.ListView;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;
import com.loopj.android.http.RequestParams;

import org.apache.http.Header;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;
import java.util.prefs.Preferences;

/**
 * Created by crashprophet on 11/30/14.
 */
public class CircleHelper {

    private Activity activity;
    private android.content.SharedPreferences prefs;
    private Resources resources;

    public CircleHelper (Activity activity, android.content.SharedPreferences prefs, Resources resources) {
        this.activity = activity;
        this.prefs = prefs;
        this.resources = resources;
    }

    public String getLocalUrlForApi () {
        AssetManager assetManager = this.resources.getAssets();
        InputStream inputStream = null;
        try {
            inputStream = assetManager.open("config.properties");
        } catch (IOException e) {
            e.printStackTrace();
        }
        Properties properties = new Properties();
        try {
            properties.load(inputStream);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return properties.getProperty("myUrl");
    }

    private class CustomAsyncHttpResponseHandler extends AsyncHttpResponseHandler {
        private JSONArray circles;

        public CustomAsyncHttpResponseHandler () {

        }

        @Override
        public void onStart () {
            System.out.println("STARTING GET REQUEST");
        }

        @Override
        public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
            String responseText = null;
            try {
                responseText = new JSONObject(new String(responseBody)).getString("results");
                this.circles = new JSONArray(responseText);
            } catch (JSONException j) {
                System.out.println(j);
            }
        }

        @Override
        public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
            String responseText = null;
            try {
                responseText = new JSONObject(new String(errorResponse)).getString("reason");
            } catch (JSONException j) {
                System.out.println(j);
            }
        }

        public JSONArray getCircles () {
            return this.circles;
        }
    }

    public JSONArray getCirclesArray () {
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(prefs);

        String userKey = "com.cherami.cherami.username";
        String username = this.prefs.getString(userKey, null);
        RequestParams params = new RequestParams();
        params.put("user", username);

        client.addHeader("Authorization", token);

        CustomAsyncHttpResponseHandler custom = new CustomAsyncHttpResponseHandler();
        client.get(this.activity.getApplicationContext(),
                   "http://" + getLocalUrlForApi() + "circles", params,
                   custom);
        return custom.getCircles();
    }
}
