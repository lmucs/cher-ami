package com.cherami.cherami;

import android.app.Activity;
import android.app.DialogFragment;
import android.app.FragmentManager;
import android.content.Context;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.net.Uri;
import android.os.Bundle;
import android.app.Fragment;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.RadioButton;
import android.widget.RadioGroup;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.io.UnsupportedEncodingException;
import java.util.Properties;

import static com.cherami.cherami.R.id.circleName;

public class CreateCircleModal extends DialogFragment {

    SharedPreferences prefs;
    Button createCircleButton;
    Button dismissModalButton;
    View root;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        Context context = getActivity().getApplicationContext();
        prefs = context.getSharedPreferences("com.cherami.cherami", Context.MODE_PRIVATE);
        super.onCreate(savedInstanceState);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {
        View rootView = inflater.inflate(R.layout.fragment_create_circle_modal, container, false);
        getDialog().setTitle("Create New Circle");
//        getDialog().setCancelable(true);

        createCircleButton = (Button) rootView.findViewById(R.id.createCircleButton);
        createCircleButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                attemptCreateCircle();
            }
        });
        dismissModalButton = (Button) rootView.findViewById(R.id.dismissModalButton);
        dismissModalButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                dismissModal();
            }
        });
        root = rootView;
        return rootView;
    }

    public JSONObject getCreateCircleParamsAsJson () {
        JSONObject jsonParams = new JSONObject();
        EditText circleName = (EditText) root.findViewById(R.id.circleName);
        RadioButton publicRadioButton = (RadioButton) root.findViewById(R.id.publicRadioButton);
        Boolean isCirclePublic = publicRadioButton.isChecked();
        System.out.println(publicRadioButton.isChecked());
        String visibilitySetting = publicRadioButton.isChecked() ? "public" : "private";

        try {
            jsonParams.put("CircleName", circleName.getText().toString());
            jsonParams.put("Public", isCirclePublic);
        } catch (JSONException j) {
            System.out.println("DONT LIKE JSON!");
        }
        System.out.println(jsonParams.toString());
        return jsonParams;
    }

    public StringEntity convertJsonUserToStringEntity (JSONObject jsonParams) {
        StringEntity entity = null;
        try {
            entity = new StringEntity(jsonParams.toString());
        } catch (UnsupportedEncodingException i) {
            System.out.println("DONT LIKE TO STRING!");
        }
        return entity;
    }

    public void dismissModal () {
        this.dismiss();
    }

    public void attemptCreateCircle() {
        AsyncHttpClient client = new AsyncHttpClient();
        String sessionKey = "com.cherami.cherami.token";
        String token = prefs.getString(sessionKey, null);

        client.addHeader("Authorization", token);
        client.post(getActivity().getApplicationContext(), ApiHelper.getLocalUrlForApi(getResources()) + "circles",
                convertJsonUserToStringEntity(getCreateCircleParamsAsJson()), "application/json",
                new AsyncHttpResponseHandler() {

                    @Override
                    public void onStart() {
                        // called before request is started
                        System.out.println("STARTING POST TO CIRCLES REQUEST");

                    }

                    @Override
                    public void onSuccess(int statusCode, Header[] headers, byte[] response) {
                        String s = new String(response);
                        // called when response HTTP status is "200 OK"

                        String responseText = null;
                        try {
                            responseText = new JSONObject(new String(response)).getString("response");
                            System.out.println(new JSONObject(new String(response)).toString());
                        } catch (JSONException j) {
                            System.out.println("Dont like JSON");
                        }

                        Log.d("Status Code: ", Integer.toString(statusCode));

                        dismissModal();
                        Toast toast = Toast.makeText(getActivity().getApplicationContext(), responseText, Toast.LENGTH_LONG);
                        toast.show();
                    }

                    @Override
                    public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                        // called when response HTTP status is "4XX" (eg. 401, 403, 404)

                        String responseText = null;
                        try {
                            responseText = new JSONObject(new String(errorResponse)).getString("Reason");
                        } catch (JSONException j) {
                            System.out.println("Dont like JSON");
                        }

                        Toast toast = Toast.makeText(getActivity().getApplicationContext(), responseText, Toast.LENGTH_LONG);
                        toast.show();
                        e.printStackTrace();
                    }

                    @Override
                    public void onRetry(int retryNo) {
                        // called when request is retried
                        System.out.println("RETRYING?!?!");
                    }
                });
    }


}