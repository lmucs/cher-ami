package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.CheckBox;
import android.widget.TextView;
import android.widget.Toast;

import org.json.JSONObject;

/**
 * Created by Geoff on 11/22/2014.
 */
public class CircleForMessageAdapter extends ArrayAdapter<CircleForMessagesItem> {

    Context context;
    int layoutResourceId;
    CircleForMessagesItem[] data = null;

    public CircleForMessageAdapter(Context context, int layoutResourceId, CircleForMessagesItem[] data) {
        super(context, layoutResourceId, data);
        this.layoutResourceId = layoutResourceId;
        this.context = context;
        this.data = data;
    }


    @Override
    public View getView(final int position, View convertView, ViewGroup parent) {
        View row = convertView;
        CircleForMessageHolder holder = null;

        if (row == null) {
            LayoutInflater inflater = ((Activity) context).getLayoutInflater();
            row = inflater.inflate(layoutResourceId, parent, false);

            holder = new CircleForMessageHolder();
            holder.txtTitle = (TextView) row.findViewById(R.id.txtTitle);
            holder.check = (CheckBox) row.findViewById(R.id.circleCheckBox);

            row.setTag(holder);

            holder.check.setOnClickListener( new View.OnClickListener() {
                public void onClick(View v) {
                    CheckBox cb = (CheckBox) v ;
                    CircleForMessagesItem circle = data[position];
                    circle.setSelected(cb.isChecked());

                    Toast.makeText(context,
                            "Clicked on Checkbox: " + cb.getText() +
                                    " is " + cb.isChecked(),
                            Toast.LENGTH_LONG).show();
                    circle.setSelected(cb.isChecked());
                }

            });
        } else {
            holder = (CircleForMessageHolder) row.getTag();
        }

        CircleForMessagesItem circleForMessage = data[position];
        try {
            holder.txtTitle.setText(circleForMessage.circleName.getString("name"));
        } catch (Exception e){
            System.out.print(e);
        }


        return row;
    }

    static class CircleForMessageHolder {
        TextView txtTitle;
        CheckBox check;
    }
}
