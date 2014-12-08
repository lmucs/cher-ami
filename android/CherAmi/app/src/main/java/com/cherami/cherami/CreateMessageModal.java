package com.cherami.cherami;

import android.app.DialogFragment;
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

    Button createMessageButton;
    Button dismissModalButton;
    View root;

    @Override
    public void onCreate(Bundle savedInstanceState) {
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
        EditText imageUrl = (EditText) root.findViewById(R.id.imageUrl);
        Bundle args = new Bundle();
        String message = messageContent.getText().toString();
        if(!imageUrl.getText().toString().equals("")) {
            message = message + " " + imageUrl.getText().toString();
        }
        args.putString("messageValue", message);
        CircleForMessageModal newFragment = new CircleForMessageModal();
        newFragment.setArguments(args);
        newFragment.show(getFragmentManager(), "dialog");
    }
}
