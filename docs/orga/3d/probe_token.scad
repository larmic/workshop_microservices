// Probe token for Story 4 (role play "Circuit Breaker").
// In HALF_OPEN only the guest holding this token may probe Flight.
// Front: raised "PROBE", back: debossed "HALF_OPEN" (mirrored, reads right
// when the token is flipped). Tab with a hole for a lanyard.
//
// Print flat, front side up. No supports needed. The raised text can get a
// second colour via a filament change at the layer where the text starts
// (z = disc_t), the AMS is not required.
// Tested with OpenSCAD 2021.01.

/* [Disc] */
disc_d    = 60;
disc_t    = 4;          // body thickness
rim_w     = 2.5;        // raised rim on the front
rim_h     = 1;

/* [Lanyard tab] */
tab_w     = 20;
tab_len   = 10;         // how far the tab sticks out of the disc
hole_d    = 6;

/* [Text] */
front_text   = "PROBE";
front_size   = 11;
front_h      = 1;       // raised
front_sub    = "NUR EINER";   // small line under the main word, "" for none
front_sub_size = 5.5;
back_text    = "HALF_OPEN";
back_size    = 6.5;
back_depth   = 0.6;     // debossed into the bottom face
font         = "DejaVu Sans:style=Bold";

$fn = 120;

hole_y = disc_d / 2 + tab_len - tab_w / 2;
echo(str("token ", disc_d, " x ", disc_d / 2 + hole_y + tab_w / 2, " x ", disc_t + max(rim_h, front_h), " mm"));

module outline() {
    hull() {
        circle(d = disc_d);
        translate([0, hole_y]) circle(d = tab_w);
    }
}

difference() {
    union() {
        linear_extrude(disc_t) outline();
        // rim on the front
        translate([0, 0, disc_t]) linear_extrude(rim_h)
            difference() { outline(); offset(delta = -rim_w) outline(); }
        // raised front text
        translate([0, 0, disc_t]) linear_extrude(front_h) {
            translate([0, front_sub == "" ? 0 : 3])
                text(front_text, size = front_size, font = font, halign = "center", valign = "center");
            if (front_sub != "")
                translate([0, -12])
                    text(front_sub, size = front_sub_size, font = font, halign = "center", valign = "center");
        }
    }
    // lanyard hole
    translate([0, hole_y, -1]) cylinder(d = hole_d, h = disc_t + rim_h + 2);
    // debossed, mirrored back text
    translate([0, 0, -1]) linear_extrude(back_depth + 1)
        mirror([1, 0, 0])
            text(back_text, size = back_size, font = font, halign = "center", valign = "center");
}
