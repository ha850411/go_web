function initProductTable()
{
    productTable = $('#productTable').on( 'processing.dt', function ( e, settings, processing ) {
        if(processing) {
            $('#loadingTable').css( 'display', 'table');
            $('#productTable_wrapper').css( 'display', 'none');
        } else {
            $('#productTable_wrapper').css( 'display', 'block');
            $('#loadingTable').css( 'display', 'none');
        }
    }).DataTable({
        language: {
            url: "static/vendor/dataTables/languages/ch-tw.json"  
        },
        serverSide:true,
        ajax: {
            url: "/api/products/getList",
            type: "GET",
        },
        columns:[
            {
                data: 'name', 
                title: "產品名稱",
            },
            {
                data: null,
                title: "數量",
                render: function(data, type, row) {
                    output = '<div style="display:flex; align-items:center;">';
                    output += '<input type="number" class="amountInput form-control" style="width:100px;" data-id="'+row.id+'" value="'+row.amount+'">';
                    output += '<button type="button" class="btn btn-success btn-sm" style="margin-left:5px" onclick="updateAmount('+row.id+')">修改</button>';
                    output += '</div>';
                    return output;
                }
            },
            {
                data: 'amountNotice', 
                title: "存貨警戒線",
                render: function(data, type, row) {
                    if(data>0) {
                        return data;
                    } else {
                        return "無";
                    }
                }
            },
            {
                data: 'formatTime',
                title: "更新時間"
            },
            {
                data: 'id', 
                title: "",
                render: function(data, type, row) {
                    output = '<div style="display:flex; align-items:center;">';
                    output += '<button type="button" class="btn btn-primary btn-sm" style="margin-right:5px">編輯產品</button>';
                    output += '<button type="button" class="btn btn-danger btn-sm">刪除產品</button>';
                    output += '</div>';
                    return output;
                }
            },
        ],
        columnDefs:[
            {
                targets: 4,
                orderable: false,
                responsivePriority: 2,
            },
            {
                targets: 0,
                responsivePriority: 1,
            }
        ],
    });
    return initProductTable;
}