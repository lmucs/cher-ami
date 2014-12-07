package com.cherami.cherami;

import android.app.DialogFragment;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.ListView;
import android.widget.Toast;
import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;
import com.loopj.android.http.RequestParams;
import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;
import java.io.UnsupportedEncodingException;

/**
 * Created by Geoff on 11/22/2014.
 */
public class CircleForMessageModal extends DialogFragment {

    CircleForMessagesItem circle_data[];
    Context context;
    Button createCircleButton;
    private ListView circleList;
    String messageValue;
    View root;
    ProgressDialog dialog;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        Bundle mArgs = getArguments();
        messageValue = mArgs.getString("messageValue");
        this.context = getActivity().getApplicationContext();
        super.onCreate(savedInstanceState);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {
        View rootView = inflater.inflate(R.layout.fragment_circle_to_post_msg_modal, container, false);
        getDialog().setTitle("Create New Message");

        getCircles(rootView);
        createCircleButton = (Button) rootView.findViewById(R.id.postButton);
        createCircleButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                boolean postMessage = false;
                for(int i = 0; i < circle_data.length; i++){
                    if(circle_data[i].isSelected()){
                        postMessage = true;
                    }
                }
                if(postMessage) {
                    attemptPostMessage();
                } else {
                    Toast toast = Toast.makeText(CircleForMessageModal.this.context, "Please select a circle to post a message to.", Toast.LENGTH_LONG);
                    toast.show();
                }
            }
        });
        root = rootView;
        return rootView;
    }

    public void getCircles(final View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(context);
        String username = ApiHelper.getUsername(context);

        RequestParams params = new RequestParams();
        params.put("user", username);

        client.addHeader("Authorization", token);
        client.get(context, ApiHelper.getLocalUrlForApi(getResources()) + "circles", params, new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {

            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                String responseText = null;
                try {
                    responseText = new JSONObject(new String(responseBody)).getString("results");
                    JSONArray y = new JSONArray(responseText);
                    circle_data = new CircleForMessagesItem[y.length()];
                    for (int x = 0; x < y.length(); x++) {
                        circle_data[x] = new CircleForMessagesItem(new JSONObject(y.get(x).toString()), false);
                    }

                    CircleForMessageAdapter adapter = new CircleForMessageAdapter(getActivity(),
                            R.layout.circle_to_post_msg, circle_data);


                    circleList = (ListView) view.findViewById(R.id.cir_msg_List);

                    circleList.setAdapter(adapter);
                } catch (JSONException j) {

                }
            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
                String responseText = null;
                try {
                    responseText = new JSONObject(new String(errorResponse)).getString("reason");
                } catch (JSONException j) {

                }

            }
        });
    }

    public JSONObject getMessageObjectRequestAsJson() {
        JSONObject jsonParams = new JSONObject();
        JSONArray circleIds = new JSONArray();
        for (int i = 0; i < circle_data.length; i++) {
            if (circle_data[i].isSelected()) {
                try {
                    circleIds.put(circle_data[i].circleName.getString("url").substring(circle_data[i].circleName.getString("url").lastIndexOf('/') + 1));
                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        }
        try {
            jsonParams.put("content", messageValue);
            jsonParams.put("circles", circleIds);
        } catch (JSONException j) {

        }

        return jsonParams;
    }

    public StringEntity convertJsonUserToStringEntity(JSONObject jsonParams) {
        StringEntity entity = null;
        try {
            entity = new StringEntity(jsonParams.toString());
        } catch (UnsupportedEncodingException i) {
            System.out.println(i);
        }

        return entity;
    }

    public void dismissModal() {
        this.dismiss();
    }

    public void attemptPostMessage() {
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(context);

        client.addHeader("Authorization", token);
        client.post(context, ApiHelper.getLocalUrlForApi(getResources()) + "messages",
                convertJsonUserToStringEntity(getMessageObjectRequestAsJson()), "application/json",
                new AsyncHttpResponseHandler() {

                    @Override
                    public void onStart() {
                        dialog = ProgressDialog.show(getActivity(), "",
                                "Loading. Please wait...", true);
                    }

                    @Override
                    public void onSuccess(int statusCode, Header[] headers, byte[] response) {
                        String responseText = null;
                        try {
                            responseText = new JSONObject(new String(response)).getString("response");
                        } catch (JSONException j) {

                        }

                        Toast toast = Toast.makeText(CircleForMessageModal.this.context, responseText, Toast.LENGTH_LONG);
                        toast.show();
                        dismissModal();
                        dialog.dismiss();
                    }

                    @Override
                    public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                        dialog.dismiss();
                        String responseText = null;

                        try {
                            responseText = new JSONObject(new String(errorResponse)).getString("reason");

                        } catch (JSONException j) {

                        }

                        Toast toast = Toast.makeText(CircleForMessageModal.this.context, responseText, Toast.LENGTH_LONG);
                        toast.show();
                        e.printStackTrace();
                        dismissModal();
                    }

                    @Override
                    public void onRetry(int retryNo) {

                    }
                });
    }
}
