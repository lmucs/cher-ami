package com.cherami.cherami;

import android.content.Context;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.content.res.Resources;

import java.io.IOException;
import java.util.Properties;

public class ApiHelper {
    private static final String cheramiPackageURL = "com.cherami.cherami";
    private static final String sessionKey = cheramiPackageURL + ".token";
    private static final String userKey = cheramiPackageURL + ".username";

    public static String getLocalUrlForApi (Resources resources) {
        AssetManager assetManager = resources.getAssets();
        Properties properties = new Properties();
        String urlPropertyKey = "myUrl";

        try {
            String propertiesFile = "config.properties";
            properties.load(assetManager.open(propertiesFile));
        } catch (IOException e) {
            e.printStackTrace();
        }

        return properties.getProperty(urlPropertyKey);
    }

    public static String getSessionToken (Context context) {
        return getSharedPreferences(context).getString(sessionKey, null);
    }

    public static String getUsername (Context context) {
        return getSharedPreferences(context).getString(userKey, null);
    }

    private static SharedPreferences getSharedPreferences (Context context) {
        return context.getSharedPreferences(cheramiPackageURL, Context.MODE_PRIVATE);
    }

    public static void saveAuthorizationToken (Context context, String userName, String token) {
        SharedPreferences prefs = getSharedPreferences(context);
        prefs.edit().putString(userKey, userName).apply();
        prefs.edit().putString(sessionKey, token).apply();
    }
}
