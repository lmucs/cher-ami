package com.cherami.cherami;

import android.app.DialogFragment;
import android.content.Context;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.EditText;

/**
 * Created by Geoff on 11/20/2014.
 */
public class CreateMessageModal extends DialogFragment {

    SharedPreferences prefs;
    Button createMessageButton;
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
        View rootView = inflater.inflate(R.layout.fragment_create_message_modal, container, false);
        getDialog().setTitle("Create New Message");

        createMessageButton = (Button) rootView.findViewById(R.id.createMessageButton);
        createMessageButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                attemptCreateMessage();
                dismissModal();
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

    public void dismissModal() {
        this.dismiss();
    }

    public void attemptCreateMessage() {
        EditText messageContent = (EditText) root.findViewById(R.id.messageContent);
        Bundle args = new Bundle();
        args.putString("messageValue", messageContent.getText().toString());
        CircleForMessageModal newFragment = new CircleForMessageModal();
        newFragment.setArguments(args);
        newFragment.show(getFragmentManager(), "dialog");
    }
}
