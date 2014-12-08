package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.graphics.Color;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.CheckBox;
import android.widget.TextView;

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
                    circle.setSelected(cb.isChecked());
                }

            });
        } else {
            holder = (CircleForMessageHolder) row.getTag();
        }

        CircleForMessagesItem circleForMessage = data[position];
        try {
            if (circleForMessage.circleName.getString("visibility").equals("private")){
                row.setBackgroundColor(Color.parseColor("#4cc1f0"));
                holder.txtTitle.setTextColor(Color.parseColor("#ffffff"));
            } else {
                row.setBackgroundColor(Color.parseColor("#f0f0f0f0"));
                holder.txtTitle.setTextColor(Color.parseColor("#000000"));
            }
            holder.txtTitle.setText(circleForMessage.circleName.getString("owner")+"'s "+circleForMessage.circleName.getString("name"));
        } catch (Exception e){

        }


        return row;
    }

    static class CircleForMessageHolder {
        TextView txtTitle;
        CheckBox check;
    }
}
