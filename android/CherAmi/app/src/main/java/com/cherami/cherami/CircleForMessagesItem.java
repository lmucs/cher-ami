package com.cherami.cherami;

import org.json.JSONObject;

/**
 * Created by Geoff on 11/22/2014.
 */
public class CircleForMessagesItem{
public JSONObject circleName;
public boolean selected;
        public CircleForMessagesItem(){
            super();
        }

        public CircleForMessagesItem(JSONObject circleName, boolean selected) {
            super();
            this.circleName = circleName;;
            this.selected = selected;
        }
        public boolean isSelected() {
            return selected;
        }
        public void setSelected(boolean selected) {
            this.selected = selected;
        }
}
