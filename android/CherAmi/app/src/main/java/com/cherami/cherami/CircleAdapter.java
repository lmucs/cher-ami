package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

public class CircleAdapter extends ArrayAdapter<CircleItem> {

    Context context;
    int layoutResourceId;
    CircleItem data[] = null;

    public CircleAdapter(Context context, int layoutResourceId, CircleItem[] data) {
        super(context, layoutResourceId, data);
        this.layoutResourceId = layoutResourceId;
        this.context = context;
        this.data = data;
    }


    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View row = convertView;
        CircleHolder holder = null;

        if (row == null) {
            LayoutInflater inflater = ((Activity) context).getLayoutInflater();
            row = inflater.inflate(layoutResourceId, parent, false);

            holder = new CircleHolder();
            holder.txtTitle = (TextView) row.findViewById(R.id.txtTitle);

            row.setTag(holder);
        } else {
            holder = (CircleHolder) row.getTag();
        }

        CircleItem circleItem = data[position];
        holder.txtTitle.setText(circleItem.title);

        return row;
    }

    static class CircleHolder {
        TextView txtTitle;
    }
}