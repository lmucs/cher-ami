package com.cherami.cherami;

import android.app.Activity;
import android.app.AlertDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;

import static android.support.v4.app.ActivityCompat.startActivity;

public class ErrorHandle {

    public static Boolean isNetworkConnected (byte[] response) {
        if (response == null) {
            return false;
        } else {
            return true;
        }
    }

    public static void displayNetworkErrorModal (Activity activity) {
        new AlertDialog.Builder(activity)
            .setTitle("Network Error")
            .setMessage("You're not connected to the network :(")
            .setNegativeButton(android.R.string.ok, new DialogInterface.OnClickListener() {
                public void onClick(DialogInterface dialog, int which) {
                    // do nothing
                }
            })
            .setIcon(android.R.drawable.ic_dialog_alert)
            .show();
    }

    public static Boolean isTokenExpired (String error) {
        if (error.equals("Missing, illegal or expired token") || error.equals("Cannot invalidate token because it is missing")) {
            return true;
        } else {
            return false;
        }
    }

    public static void displayTokenModal (final Activity activity) {
        new AlertDialog.Builder(activity)
                .setTitle("Session Expired")
                .setMessage("You're session has expired. Please login again.")
                .setNegativeButton(android.R.string.ok, new DialogInterface.OnClickListener() {
                    public void onClick(DialogInterface dialog, int which) {
                        // Redirect to login
                        Intent intent = new Intent(activity.getApplicationContext(), LoginActivity.class);
                        activity.startActivity(intent);
                        activity.finish();
                    }
                })
                .setIcon(android.R.drawable.ic_dialog_alert)
                .show();
    }
}
