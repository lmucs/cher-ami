package com.cherami.cherami;


import android.app.DialogFragment;
import android.app.ProgressDialog;
import android.content.Context;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.EditText;
import android.widget.RadioButton;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONException;
import org.json.JSONObject;
import java.io.UnsupportedEncodingException;

public class CreateCircleModal extends DialogFragment {

    Context context;
    Button createCircleButton;
    Button dismissModalButton;
    View root;
    ProgressDialog dialog;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        this.context = getActivity().getApplicationContext();
        super.onCreate(savedInstanceState);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {
        View rootView = inflater.inflate(R.layout.fragment_create_circle_modal, container, false);
        getDialog().setTitle("Create New Circle");

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
                CreateCircleModal.this.dismiss();
            }
        });
        root = rootView;

        return rootView;
    }

    public String getCircleName () {
        return ((EditText) root.findViewById(R.id.circleName)).getText().toString();
    }

    public JSONObject getCreateCircleParamsAsJson () {
        JSONObject jsonParams = new JSONObject();
        RadioButton publicRadioButton = (RadioButton) root.findViewById(R.id.publicRadioButton);
        Boolean isCirclePublic = publicRadioButton.isChecked();
        String visibilitySetting = publicRadioButton.isChecked() ? "public" : "private";

        try {
            jsonParams.put("CircleName", getCircleName());
            jsonParams.put("Public", isCirclePublic);
        } catch (JSONException j) {

        }

        return jsonParams;
    }

    public StringEntity convertJsonUserToStringEntity (JSONObject jsonParams) {
        StringEntity entity = null;

        try {
            entity = new StringEntity(jsonParams.toString());
        } catch (UnsupportedEncodingException i) {

        }

        return entity;
    }

    public void attemptCreateCircle() {
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(context);

        client.addHeader("Authorization", token);
        client.post(context, ApiHelper.getLocalUrlForApi(getResources()) + "circles",
                convertJsonUserToStringEntity(getCreateCircleParamsAsJson()), "application/json",
                new AsyncHttpResponseHandler() {

                    @Override
                    public void onStart() {
                        dialog = ProgressDialog.show(getActivity(), "",
                                "Loading. Please wait...", true);
                    }

                    @Override
                    public void onSuccess(int statusCode, Header[] headers, byte[] response) {
                        String successfulCreation = "Successfully created new " +
                                                    CreateCircleModal.this.getCircleName() +
                                                    " circle.";
                        CreateCircleModal.this.dismiss();
                        Toast toast = Toast.makeText(getActivity().getApplicationContext(),
                                                     successfulCreation, Toast.LENGTH_LONG);
                        toast.show();
                        dialog.dismiss();

                    }

                    @Override
                    public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable e) {
                        dialog.dismiss();
                        String responseText = null;

                        try {
                            if (!ErrorHandle.isNetworkConnected(errorResponse)) {
                                ErrorHandle.displayNetworkErrorModal(getActivity());

                            } else {
                                responseText = new JSONObject(new String(errorResponse)).getString("reason");
                                if (ErrorHandle.isTokenExpired(responseText)) {
                                    ErrorHandle.displayTokenModal(getActivity());
                                }
                            }
                        } catch (JSONException j) {

                        }
                    }

                    @Override
                    public void onRetry(int retryNo) {

                    }
                });
    }
}