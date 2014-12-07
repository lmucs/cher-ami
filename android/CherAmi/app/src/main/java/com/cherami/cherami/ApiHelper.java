package com.cherami.cherami;

import android.content.Context;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.content.res.Resources;

import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

/**
 * Created by crashprophet on 11/30/14.
 */
public class ApiHelper {



    public static String getLocalUrlForApi (Resources resources) {
        AssetManager assetManager = resources.getAssets();
        Properties properties = new Properties();
        try {
            properties.load(assetManager.open("config.properties"));
        } catch (IOException e) {
            e.printStackTrace();
        }

        return properties.getProperty("myUrl");
    }

    public static String getSessionToken (Context context) {
        String sessionKey = "com.cherami.cherami.token";
        return getSharedPreferences(context).getString(sessionKey, null);
    }

    public static String getUsername (Context context) {
        String usernameKey = "com.cherami.cherami.username";
        return getSharedPreferences(context).getString(usernameKey, null);
    }

    private static SharedPreferences getSharedPreferences (Context context) {
        return context.getSharedPreferences("com.cherami.cherami", Context.MODE_PRIVATE);
    }

    public static void saveAuthorizationToken (Context context, String userName, String token) {
        String sessionKey = "com.cherami.cherami.token";
        String userKey = "com.cherami.cherami.username";
        SharedPreferences prefs = getSharedPreferences(context);

        prefs.edit().putString(userKey, userName).apply();
        prefs.edit().putString(sessionKey, token).apply();
    }
}
