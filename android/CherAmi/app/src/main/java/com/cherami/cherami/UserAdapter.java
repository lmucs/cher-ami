package com.cherami.cherami;

import android.app.Activity;
import android.content.ClipData;
import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;

import org.json.JSONException;

public class UserAdapter extends ArrayAdapter<User> {

    Context context;
    int layoutResourceId;
    User[] data = null;

    public UserAdapter(Context context, int layoutResourceId, User[] data) {
        super(context, layoutResourceId, data);
        this.layoutResourceId = layoutResourceId;
        this.context = context;
        this.data = data;
    }


    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View row = convertView;
        UserHolder holder = null;

        if (row == null) {
            LayoutInflater inflater = ((Activity) context).getLayoutInflater();
            row = inflater.inflate(layoutResourceId, parent, false);

            holder = new UserHolder();
            holder.txtTitle = (TextView) row.findViewById(R.id.txtTitle);

            row.setTag(holder);
        } else {
            holder = (UserHolder) row.getTag();
        }

        User user = data[position];

        try {
            holder.txtTitle.setText(user.userName.getString("u.handle"));
        } catch (JSONException e) {
            e.printStackTrace();
        }

        return row;
    }

    @Override
    public User getItem(int position){
     return data[position];
    }

    static class UserHolder {
        TextView txtTitle;
    }
}
